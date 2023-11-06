package builder

import (
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB

type Foods struct {
    FdcId               int             `json:"fdc_id" gorm:"primary_key"`
    Description         string          `json:"description"`
    Calories            float32         `json:"calories"`
    CaloriesFromFat     float32         `json:"calories_from_fat"`
    TotalFat            float32         `json:"total_fat"`
    TransFat            float32         `json:"trans_fat"`
    SaturatedFat        float32         `json:"saturated_fat"`
    Cholesterol         float32         `json:"cholesterol"`
    Sodium              float32         `json:"sodium"`
    Potassium           float32         `json:"potassium"`
    TotalCarbs          float32         `json:"total_carbs"`
    DietaryFiber        float32         `json:"dietary_fiber"`
    Sugars              float32         `json:"sugars"`
    Protein             float32         `json:"protein"`
    VitaminA            float32         `json:"vitamin_a"`
    VitaminB6           float32         `json:"vitamin_b6"`
    VitaminB12          float32         `json:"vitamin_b12"`
    VitaminC            float32         `json:"vitamin_c"`
    VitaminD            float32         `json:"vitamin_d"`
    VitaminE            float32         `json:"vitamin_e"`
    Biotin              float32         `json:"biotin"`
    Niacin              float32         `json:"niacin"`
    Riboflavin          float32         `json:"riboflavin"`
    Thiamin             float32         `json:"thiamin"`
    Copper              float32         `json:"copper"`
    Calcium             float32         `json:"calcium"`
    Iron                float32         `json:"iron"`
    Magnesium           float32         `json:"magnesium"`
    Phosphorus          float32         `json:"phosphorus"`
    Iodine              float32         `json:"iodine"`
    Zinc                float32         `json:"zinc"`
}

type Portion struct {
    Pid             int         `json:"pid" gorm:"primary_key"`
    FdcId           int         `json:"fdc_id"` 
    Amount          float32     `json:"amount"`
    UnitName        string      `json:"unit_name"`
    AbbrUnitName    string      `json:"abbr_unit_name"`
    GramWeight      float32     `json:"gram_weight"`
}

func ConfigureDB() error {
    dsn := os.Getenv("DSN")
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if (err != nil) {
        return err 
    } 
    
    Db = db;
    return nil
}

func GetFood(food string) Foods {
    var usdaFood Foods 
    wildcard := "%" + food + "%"
    Db.Where("description LIKE ?", wildcard).First(&usdaFood)
    return usdaFood
}

func GetAvailablePortions(fdcId int) []Portion {
    var portions []Portion
    Db.Raw("select * from portion WHERE fdc_id = ?", fdcId).Scan(&portions)
    return portions
}
