package parser

import (
    "os"
    "bufio"
    "strings"
    "log"
    "math/rand"
    "golang.org/x/text/cases"
    "golang.org/x/text/language"
    "regexp"
    "strconv"
)

const QTY_REGEX = `[0-9]+\.*[0-9]*`

func LoadUagents(uagents *[]string) error {
    f, err := os.Open("/usr/src/backend/uagents.txt")
    if err != nil {
        return err
    }
    
    // We need to read the entire file into an array because there
    // is no mechanism of getting the total lines with the File API
    // like in Python.
    scanner := bufio.NewScanner(f)
    var n int = 10000
    for i := 0; i < n; i++ {
        if ok := scanner.Scan(); ok {
            *uagents = append(*uagents, scanner.Text())
        }
    }

    return nil
}

func GetUnit(key string) string {
    if key == "calories" {
        return "cals"
    } else if key == "carbohydrateContent" ||
    key == "fatContent" || key == "fiberContent" ||
    key == "proteinContent" || key == "sugarContent" ||
    key == "transFatContent" || key == "unsaturatedFatContent" ||
    key == "saturatedFatContent" {
        return "g"
    } else {
        return "mg"
    }
}

func CreateName(key string) string {
    c := cases.Title(language.English)

    if key == "saturatedFatContent" {
        return "Saturated fat"
    } else if key == "transFatContent" {
        return "Trans fat"
    } else if key == "unsaturatedFatContent" {
        return "Unsaturated fat"
    } else if strings.HasSuffix(key, "Size") {
        return c.String(strings.Replace(key, "Size", "", 1))
    } else if strings.HasSuffix(key, "Content") {
        return c.String(strings.Replace(key, "Content", "", 1))
    } else {
        return c.String(key)
    }
}

func FindRecipe(json interface{}) map[string]interface{} {
    // REQUIRES:    json
    // MODIFIES:    none
    // EFFECTS:     Checks the type of the json (using type coercion) and then
    //              gets the recipe schema based on that type

    switch json := json.(type) {
    case map[string]interface{}:
        return findInMap(json)
    case []interface{}:
        return findInInterfaceList(json)
    case []map[string]interface{}:
        return findInMapList(json)
    default:
        log.Printf("[PARSER] recipe error: encountered unexpected type %T\n", json)
        return nil
    }
}

func isRecipeSchema (schemaType interface{}) bool {
    // REQUIRES:    schemaType
    // MODIFIES:    none
    // EFFECTS:     if the schemaType is a string, it will check if it is a
    //              recipe. If it's an array, find "recipe" inside the array
    
    switch t := schemaType.(type) {
    case string:
        schemaType = strings.ToLower(t)
        if schemaType == "recipe" {
            return true 
        }
    case []interface{}:
        for _, val := range t {
            if val, ok := val.(string); ok {
                val = strings.ToLower(val)
                if val == "recipe" {
                    return true
                }
            }
        }
    default:
        log.Printf("[PARSER] recipe normalization: encountered unexpected type %T\n", schemaType)
    }
    
    return false
} 

func findInMap(json map[string]interface{}) map[string]interface{} {
    // REQUIRES:    json
    // MODIFIES:    none
    // EFFECTS:     When the JSON is a map, it will either be a map containing 
    //              a list or the regular map. If it contains a list, we should
    //              loop through the list and find the desired schema. If it's
    //              just a regular map, we should make sure the map is a holding
    //              recipe data

    if nodeArray, exists := json["@graph"]; exists {
        return FindRecipe(nodeArray) 
    }
    
    if schemaType, exists := json["@type"]; exists {
        if isRecipeSchema(schemaType) {
            return json
        }
    }
    
    return nil
}

func findInInterfaceList(json []interface{}) map[string]interface{} {
    // REQUIRES:    json 
    // MODIFIES:    none
    // EFFECTS:     When the JSON is a map, it will either be a map containing 
    //              a list or the regular map. If it contains a list, we should
    //              loop through the list and find the desired schema. If it's
    //              just a regular map, we should make sure the map is a holding
    //              recipe data

    for _, schema := range json {
        if schema, ok := schema.(map[string]interface{}); ok {
            if schemaType, exists := schema["@type"]; exists {
                if isRecipeSchema(schemaType) {
                    return schema 
                }
            }
        }
    }
    
    return nil
}

func findInMapList(json []map[string]interface{}) map[string]interface{} {
    // REQUIRES:    json 
    // MODIFIES:    none
    // EFFECTS:     When the JSON is a map, it will either be a map containing 
    //              a list or the regular map. If it contains a list, we should
    //              loop through the list and find the desired schema. If it's
    //              just a regular map, we should make sure the map is a holding
    //              recipe data

    for _, schema := range json {
        if schemaType, exists := schema["@type"]; exists && isRecipeSchema(schemaType) {
            if isRecipeSchema(schemaType) {
                return schema 
            }
        }
    }

    return nil
}

func BuildNutrientJSON(key, val string) map[string]interface{} {
    exp := regexp.MustCompile(QTY_REGEX)
    
    match := exp.FindIndex([]byte(val))

    if len(match) == 2 {
        i, j := match[0], match[1]
        qty, err := strconv.ParseFloat(val[i:j], 64)
        if err != nil {
            log.Printf("error: %+v", err)
        }

        unit, name := GetUnit(key), CreateName(key)
        return map[string]interface{}{
            "qty": qty,
            "unit": unit,
            "name": name,
        }
    }

    return nil
}

func NormalizeNutritionData(nutrition interface{}) interface{} {
    // REQUIRES:    nutrition
    // MODIFIES:    nutrition
    // EFFECTS:     Converts [nutrition] into a map, and modifies it's
    //              keys so that it matches the JSON format of 
    //              {"val": "", "unit": "", "name": ""}

    if nutrition, ok := nutrition.(map[string]interface{}); ok {
        delete(nutrition, "@type")
        delete(nutrition, "@context")

        for key, val := range nutrition {
            if val == nil {
                delete(nutrition, key)
            }
            
            if val, ok := val.(string); ok {
                if json := BuildNutrientJSON(key, val); json != nil {
                    nutrition[key] = json
                }
            }
        }
    }

    return nutrition
}

func NormalizeImageData(img interface{}) interface{} {
    // REQUIRES:    img 
    // MODIFIES:    img 
    // EFFECTS:     Peeks into the [img] interface and looks for a string
    //              that represents the url to an image of the recipe. Url can be
    //              in the form of a list, map, or string
    
    switch img := img.(type) {
    case map[string]interface{}:
        if url, exists := img["url"]; exists {
            return url
        }
    case []interface{}:
        i := rand.Intn(len(img))
        switch imgObj := img[i].(type) {
        case map[string]interface{}:
            if link, exists := imgObj["url"]; exists {
                return link
            }
        case string:
            return imgObj
        default:
            log.Printf("[PARSER] image url: encountered type %T\n", imgObj)
        }
    default:
        log.Printf("[PARSER] image normalization: encountered type %T\n", img)
    }

    return img 
}

func NormalizeMainEntity(mainEntity interface{}) interface{} {
    // REQUIRES:    mainEntity 
    // MODIFIES:    mainEntity 
    // EFFECTS:     Peeks into the [mainEntity] interface and looks for a string
    //              that represents the url to the recipe. Url can be
    //              in the form of a list, map, or string
    
    if mainEntity, ok := mainEntity.(map[string]interface{}); ok {
        if url, exists := mainEntity["@id"]; exists {
            return url
        }
    }

    return mainEntity 
}

