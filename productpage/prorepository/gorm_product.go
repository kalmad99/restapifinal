package prorepository

import (
	"../../entity"
	"../../productpage"
	"errors"
	"github.com/jinzhu/gorm"
	"log"
	"math"
)

// ItemGormRepo implements the menu.ItemRepository interface
type ItemGormRepo struct {
	conn *gorm.DB
}

// NewItemGormRepo will create a new object of ItemGormRepo
func NewItemGormRepo(db *gorm.DB) productpage.ItemRepository {
	return &ItemGormRepo{conn: db}
}

// Items returns all food menus stored in the database
func (itemRepo *ItemGormRepo) Items() ([]entity.Product, []error) {
	items := []entity.Product{}
	errs := itemRepo.conn.Find(&items).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return items, errs
}

// Item retrieves a food menu by its id from the database
func (itemRepo *ItemGormRepo) Item(id uint) (*entity.Product, []error) {
	item := entity.Product{}
	errs := itemRepo.conn.First(&item, id).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return &item, errs
}

// UpdateItem updates a given food menu item in the database
func (itemRepo *ItemGormRepo) UpdateItem(item *entity.Product) (*entity.Product, []error) {
	itm := item
	errs := itemRepo.conn.Save(itm).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return itm, errs
}

// DeleteItem deletes a given food menu item from the database
func (itemRepo *ItemGormRepo) DeleteItem(id uint) (*entity.Product, []error) {
	itm, errs := itemRepo.Item(id)

	if len(errs) > 0 {
		return nil, errs
	}

	errs = itemRepo.conn.Delete(itm, id).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return itm, errs
}

// StoreItem stores a given food menu item in the database
func (itemRepo *ItemGormRepo) StoreItem(item *entity.Product) (*entity.Product, []error) {
	itm := item
	errs := itemRepo.conn.Create(itm).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return itm, errs
}
func (itemRepo *ItemGormRepo) SearchProduct(index string) ([]entity.Product, error) {
	items := []entity.Product{}

	err := itemRepo.conn.Where("name ILIKE ?", "%"+index+"%").Find(&items).GetErrors()

	if len(err) != 0 {
		//return nil, err
		errors.New("Search Product Repo not working")
	}
	return items, nil
}
func (itemRepo *ItemGormRepo) RateProduct(pro *entity.Product) (*entity.Product, []error) {
	var oldratings float64
	var oldcount float64

	item := &entity.Product{}
	//itemRepo.conn.Where("id = ?", pro.ID).First(&item)
	row := itemRepo.conn.Select("rating").First(&item).Scan(&oldratings)
	item.Rating = oldratings
	log.Println("rating", item.Rating)
	if row.RecordNotFound(){
		panic(row.Error)
	}

	row = itemRepo.conn.Select("raters_count").First(&item).Scan(&oldcount)
	item.RatersCount = oldcount
	log.Println("count", item.RatersCount)
	if row.RecordNotFound(){
		panic(row.Error)
	}

	var newratings = ((oldratings*oldcount) + pro.Rating)/(oldcount+1)

	row = itemRepo.conn.Model(&item).UpdateColumns(entity.Product{Rating: float64(math.Round(newratings*2))/2,RatersCount: oldcount+1})

	if row.RowsAffected < 1{
		return nil, row.GetErrors()
	}
	return item, nil
}
