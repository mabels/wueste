cd smoke
if [ $? != 0 ]
then
  echo "not in smoke dir"
  exit 1
fi
rm -f package.json
pnpm init
pnpm install -f ../pubdir/wueste-*.tgz
sh node_modules/.bin/wueste-generator --write-test-schema --input-file ./generated/simple_type.schema.json --input-file ./node_modules/wueste/payload.schema.json --output-dir ./generated
sh node_modules/.bin/wueste-generate-group-type --input-files generated/simple_type.schema.json --output-dir ./generated --include-path ./generated
npm_config_yes=true npx ts-node ./smoke.ts

