cd smoke
if [ $? != 0 ]
then
  echo "not in smoke dir"
  exit 1
fi
rm -rf package.json generated
pnpm init
pnpm install -f ../pubdir/wueste-*.tgz
sh node_modules/.bin/wueste-generator \
  --write-test-schema \
  --input-file ./generated/simple_type.schema.json \
  --input-file ./node_modules/wueste/payload.schema.json \
  --output-dir ./generated
sh node_modules/.bin/wueste-generate-group-type \
  --input-files generated/simple_type.schema.json \
  --output-dir ./generated/ts \
  --include-path ./generated
sh node_modules/.bin/wueste-generate-group-type \
  --input-files generated/simple_type.schema.json \
  --output-dir ./generated/jschema \
  --output-format JSchema
sh node_modules/.bin/wueste-generator \
  --input-file ./generated/jschema/simpletypekey.schema.json \
  --output-dir ./generated/jschema
npm_config_yes=true npx ts-node ./smoke.ts

