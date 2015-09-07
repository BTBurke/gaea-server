package email

import "fmt"

const InventoryOutOfStockTemplate string = `
Hello {{.first_name}},

Unfortunately, an item you recently ordered is now out of stock.  The item you ordered is:

{{.item_name}}

Because the item is currently unavailable, we've temporarily removed it from your order.  However, if the item comes back in stock before we finalize the order with the vendor, we'll send you an email and add it back into your order.  If you would rather substitute a different item for this one, you can log on to the website and choose a replacement item as long as the sale is still open.

https://guangzhouaea.org

If you have any questions, you can email us at <a href="mailto:orders@guangzhouaea.org">orders@guangzhouaea.org</a>

Sincerely,
Guangzhou AEA
`

func InventoryOutOfStock(firstName string, itemName string) (string, error) {
	data := map[string]string{
		"first_name": firstName,
		"item_name":  itemName,
	}
	body, err := RenderFromTemplate(data, InventoryOutOfStockTemplate)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	fmt.Println(body)
	return body, nil
}
