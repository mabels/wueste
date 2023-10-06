# wueste

Wueste provides and a generator which can convert a json-schema to "typescript" with easy way to set attributes and get a Constrained Object back.

This is done with the Builder Pattern like:

```
interface Object {
    readonly Attr: string
    readonly OptionalAttr?: string
}



const obj = new ObjectBuilder().
                    Attr("value").
                    OptionalAttr("x").
                    Get()
```

this works nice on flat objects but not on nested objects.

```
interface NestedObject {
    readonly Nested: Object
}

const obj = new NestedObjectBuilder()
                    .Nested(new ObjectBuilder()
                            .Attr("value")
                            .Get()
                    ).Get()
```

or

```
const obj = new NestedObjectBuilder().
        Nested((nestedBuilder) => {
            nestedBuilder.Attr("value")
            nestedBuilder.OptionalAttr("x")
        }).Get()
```

or any better way to do this?

The generator enables the abiliy to get back every attribute in a defined order.
To provide a way create content-addressable with a hash function.

This is done with:

```
ObjectGetter(obj).Apply((path, value) => {
    console.log(path, value)
})
```

I will provide a way to set the attributes in the simlar way like the Getter works.

```
interface ObjectBuilder {
    Attr(value: string): ObjectBuilder
    OptionalAttr(value: string): ObjectBuilder
    Get(): Object
}

ObjectFactory.Builder().Setter((path, value) => {

}).Get()

ObjectSetter(builder).Apply((path, value) => {
    console.log(path, value)
    return WuestenRetVal(value)
})

/* object:NoName:
    const v0 = v
    builder.Test(fn(
      [
        ...base,
        helperTestSchema,
        (helperTestSchema as WuestenReflectionObject).properties![0],
      ]
    , builder.GetTest()))
*/
```

```




```
