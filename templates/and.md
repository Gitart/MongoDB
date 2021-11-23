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
