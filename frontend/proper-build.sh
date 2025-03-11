#! /bin/bash

echo -e "Cleaning dist..."
find dist -mindepth 1 -delete

echo -e "Running nuxi generate..."
npm run generate

echo -e "Copying data to dist..."
cp -r .output/public/. dist/