package supplylist

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"school-supply-list/api/supplies"
	"school-supply-list/database"
	"strconv"
)

type SupplyList struct {
	ListID              int                              `json:"list_id"`
	Grade               int                              `json:"grade"`
	SchoolID            int                              `json:"school_id"`
	ListName            string                           `json:"list_name"`
	BasicSupplies       []supplies.SupplyItem            `json:"basic_supplies"`
	CategorizedSupplies map[string][]supplies.SupplyItem `json:"categorized_supplies"`
	Published           bool                             `json:"published"`
}

func CreateSupplyList(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var list SupplyList
		err := c.BindJSON(&list)
		if err != nil {
			c.AbortWithStatusJSON(400, "Invalid request.")
			return
		}
		row := db.Db.QueryRow(`INSERT INTO supply_list (grade, list_name, school_id, published, list_id) 
		  VALUES ($1, $2, $3, false, default) returning list_id`, list.Grade, list.ListName, list.SchoolID)
		err = row.Scan(&list.ListID)
		if err != nil {
			database.CheckDBErr(err.(*pq.Error), c)
			return
		}

		c.JSON(200, list)
	}
}

func GetSupplyList(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idString := c.Param("id")
		id, err := strconv.Atoi(idString)
		if err != nil {
			c.AbortWithStatusJSON(400, "Invalid id. Must be an integer.")
			return
		}
		list := SupplyList{
			ListID: id,
		}

		rowCount := -1
		rows, err := db.Db.Query(`SELECT list_id, grade, list_name, school_id from supply_list 
											where supply_list.list_id=$1`, id)
		if err != nil {
			database.CheckDBErr(err.(*pq.Error), c)
			return
		}
		for rows.Next(){
			err = rows.Scan(&list.ListID, &list.Grade, &list.ListName, &list.SchoolID)
			if err != nil {
				fmt.Println(err)
				database.CheckDBErr(err.(*pq.Error), c)
				return
			}
			rowCount ++
		}

		if rowCount == -1{
			c.AbortWithStatusJSON(404, "This resource does not exist")
			return
		}

		list.BasicSupplies, list.CategorizedSupplies, err = getItemsForList(list.ListID, db)
		if err != nil {
			database.CheckDBErr(err.(*pq.Error), c)
			return
		}

		c.JSON(200, list)
	}
}

func getItemsForList(id int, db *database.DB) ([]supplies.SupplyItem, map[string][]supplies.SupplyItem, error) {
	var basicSupplies []supplies.SupplyItem
	categorizedSupplies := make(map[string][]supplies.SupplyItem)
	rows, err := db.Db.Query(`SELECT id, supply_name, supply_desc, ilb.category FROM supply_item sup 
										INNER JOIN item_list_bridge ilb on sup.id = ilb.item_id
										WHERE ilb.list_id = $1`, id)
	if err != nil {
		return basicSupplies, categorizedSupplies, err
	}
	for rows.Next() {
		var supply supplies.SupplyItem

		err = rows.Scan(&supply.ID, &supply.Supply, &supply.Desc, &supply.Category)
		basicSupplies = append(basicSupplies, supply)

		// Check if item is categorized
		if supply.Category.Valid {
			// Either create a new map item or add to existing item
			if val, ok := categorizedSupplies[supply.Category.String]; ok {
				val = append(val, supply)
			} else {
				categorizedSupplies[supply.Category.String] = []supplies.SupplyItem{supply}
			}
		} else {
			basicSupplies = append(basicSupplies, supply)
		}
	}
	return basicSupplies, categorizedSupplies, nil
}

func UpdateSupplyList(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idString := c.Param("id")
		id, err := strconv.Atoi(idString)
		if err != nil {
			c.AbortWithStatusJSON(400, "Invalid id. Must be an integer.")
			return
		}
		var list SupplyList
		err = json.NewDecoder(c.Request.Body).Decode(&list)

		list.ListID = id

		if err != nil {
			c.AbortWithStatusJSON(400, "Invalid request.")
			return
		}
		row := db.Db.QueryRow(`UPDATE supply_list SET grade=$1, list_name=$2, published=$3, 
	   		school_id=$4 where list_id=$5 returning list_id`, list.Grade, list.ListName, list.Published, list.SchoolID, list.ListID)
		err = row.Scan(&list.ListID)
		if err != nil {
			database.CheckDBErr(err.(*pq.Error), c)
			return
		}

		c.JSON(200, list)
	}
}

func DeleteSupplyList(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idString := c.Param("id")
		id, err := strconv.Atoi(idString)
		if err != nil {
			c.AbortWithStatusJSON(400, "Invalid id. Must be an integer.")
			return
		}
		row := db.Db.QueryRow(`DELETE FROM supply_list where list_id=$1`, id)
		if row.Err() != nil {
			database.CheckDBErr(err.(*pq.Error), c)
			return
		}

		c.JSON(200, nil)
	}
}
