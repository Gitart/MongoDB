## Filter data 


```
type Sale struct {
    ProductName string    `bson:"product_name"`
    Price       int       `bson:"price"`
    SaleDate    time.Time `bson:"sale_date"`
}
```

Then you can query it like this:


```sql
fromDate := time.Date(2014, time.November, 4, 0, 0, 0, 0, time.UTC)
toDate := time.Date(2014, time.November, 5, 0, 0, 0, 0, time.UTC)

var sales_his []Sale
err = c.Find(
    bson.M{
        "sale_date": bson.M{
            "$gt": fromDate,
            "$lt": toDate,
        },
    }).All(&sales_his)
```


## filter time now

```sql
cur, err := collection.Find(ctx, bson.M{"createdAt": bson.M{
    "$gte": primitive.NewDateTimeFromTime(time.Now().AddDate(-1, 0, 0)),
}})
```
