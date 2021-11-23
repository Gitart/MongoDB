## Usege $and

```sql
ageLower, err := strconv.Atoi(filters["ageLower"])
ageUpper, err := strconv.Atoi(filters["ageUpper"])
heightLower, err := strconv.Atoi(filters["heightLower"])
heightUpper, err := strconv.Atoi(filters["heightUpper"])
if err != nil {
    return nil, err
}

dobUpper := time.Now().AddDate(-ageLower, 0, 0)
dobLower := time.Now().AddDate(-ageUpper, 0, 0)

pColl := s.DB("mydb").C("profiles")

query := bson.M{
    "$and": []bson.M{
        bson.M{"active": bson.M{"$eq": true}},
        bson.M{"gender": bson.M{"$ne": u.Gender}},
        bson.M{"_id": bson.M{"$nin": u.HiddenProfiles}},
        bson.M{"_id": bson.M{"$ne": u.ProfileID}},
        bson.M{"dob": bson.M{"$gt": dobLower , "$lt": dobUpper}},
        bson.M{"height": bson.M{"$gt": heightLower, "$lt": heightUpper}},
    },
}

return pColl.Find(query).Select(bson.M{"first_name": 0}).All(&profiles)ageLower, err := strconv.Atoi(filters["ageLower"])
ageUpper, err := strconv.Atoi(filters["ageUpper"])
heightLower, err := strconv.Atoi(filters["heightLower"])
heightUpper, err := strconv.Atoi(filters["heightUpper"])
if err != nil {
    return nil, err
}

dobUpper := time.Now().AddDate(-ageLower, 0, 0)
dobLower := time.Now().AddDate(-ageUpper, 0, 0)

pColl := s.DB("mydb").C("profiles")

query := bson.M{
    "$and": []bson.M{
        bson.M{"active": bson.M{"$eq": true}},
        bson.M{"gender": bson.M{"$ne": u.Gender}},
        bson.M{"_id": bson.M{"$nin": u.HiddenProfiles}},
        bson.M{"_id": bson.M{"$ne": u.ProfileID}},
        bson.M{"dob": bson.M{"$gt": dobLower , "$lt": dobUpper}},
        bson.M{"height": bson.M{"$gt": heightLower, "$lt": heightUpper}},
    },
}

return pColl.Find(query).Select(bson.M{"first_name": 0}).All(&profiles)
```


## Iso date

```sql
collection := client.Database("Hello").Collection("demo")
    eventStartTime := "2013-10-01T01:11:18.965Z" //string format
    eventEndTime := "2014-12-03T01:11:18.965Z"   //string format

    const (
        layoutISO = "2006-01-02T15:04:05.000Z"
    )

    //import "time" package.
    t1, _ := time.Parse(layoutISO, eventStartTime) //converted to ISODate format
    t2, _ := time.Parse(layoutISO, eventEndTime)   //converted to ISODate format
    //fmt.Println(t1)
    filterCursor, err := collection.Find(context.TODO(), bson.M{"sentat": bson.M{"$gt": t1, "$lt": t2}})
    if err != nil {
        log.Fatal(err)
    }
    var result []bson.M
    if err = filterCursor.All(context.TODO(), &result); err != nil {
        log.Fatal(err)
    }
    fmt.Println(result)
    ```
    
    
