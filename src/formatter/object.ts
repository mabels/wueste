export class WuestenAttributeObject<T, C, O> extends WuestenAttr<T, C, O> {
  private readonly _builder: WuestenAttribute<T, C>;
  constructor(param: WuestenAttributeParameter<C>, factory: WuestenFactory<T, C, O>) {
    const builder = factory.Builder(param);
    super(param, { coerce: builder.Coerce.bind(builder) });
    this._builder = builder;
  }

  Coerce(value: C): Result<T> {
    const res = this._builder.Coerce(value);
    if (res.is_ok()) {
      this._value = res.unwrap();
    }
    return res;
  }

  Get(): Result<T> {
    if (this.param.default === undefined && this._value === undefined) {
      return Result.Err(`Attribute[${WuestenAttributeName(this.param)}] is required`);
    }
    if (this._value !== undefined) {
      return Result.Ok(this._value);
    }
    // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
    return Result.Ok(this.param.default! as T);
  }
}
