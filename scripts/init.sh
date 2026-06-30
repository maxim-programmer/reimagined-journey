#!/bin/bash

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${GREEN}📚 Initializing system with test documents...${NC}"

echo -e "${YELLOW}🔐 Registering user...${NC}"
curl -s -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"login":"testuser","password":"test123"}' > /dev/null 2>&1 || true

echo -e "${YELLOW}🔐 Logging in...${NC}"
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"login":"testuser","password":"test123"}')

TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo -e "${RED}❌ Failed to get token${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Token obtained${NC}"

mkdir -p ./temp_docs

echo -e "${YELLOW}📥 Downloading 10 test PDF files...${NC}"

urls=(
    "https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf"
    "https://www.africau.edu/images/default/sample.pdf"
    "https://www.irs.gov/pub/irs-pdf/p1544.pdf"
    "https://www.clickdimensions.com/links/TestPDFfile.pdf"
    "https://www.pdf995.com/samples/pdf.pdf"
    "https://www.orimi.com/pdf-test.pdf"
    "https://pdfobject.com/pdf/sample.pdf"
    "https://filesamples.com/samples/document/pdf/sample1.pdf"
    "https://filesamples.com/samples/document/pdf/sample2.pdf"
    "https://filesamples.com/samples/document/pdf/sample3.pdf"
)

total=${#urls[@]}
downloaded=0

for i in "${!urls[@]}"; do
    filename="doc_$((i+1)).pdf"
    echo -e "${YELLOW}Downloading $((i+1))/$total: $filename...${NC}"

    if curl -L -s --max-time 30 -o "./temp_docs/$filename" "${urls[$i]}"; then
        if [ -s "./temp_docs/$filename" ]; then
            echo -e "${GREEN}✅ Downloaded $filename ($(du -h ./temp_docs/$filename | cut -f1))${NC}"
            downloaded=$((downloaded + 1))
        else
            echo -e "${RED}❌ Empty file: $filename${NC}"
            rm -f "./temp_docs/$filename"
        fi
    else
        echo -e "${RED}❌ Failed to download: ${urls[$i]}${NC}"
        rm -f "./temp_docs/$filename"
    fi
done

echo -e "${GREEN}📤 Uploading documents...${NC}"

uploaded=0
for file in ./temp_docs/*.pdf; do
    if [ -f "$file" ] && [ -s "$file" ]; then
        echo -e "${YELLOW}Uploading $(basename "$file")...${NC}"
        response=$(curl -s -X POST http://localhost:8080/api/v1/documents/upload \
            -H "Authorization: Bearer $TOKEN" \
            -F "file=@$file")

        if echo "$response" | grep -q '"id"'; then
            echo -e "${GREEN}✅ Successfully uploaded $(basename "$file")${NC}"
            uploaded=$((uploaded + 1))
        else
            echo -e "${RED}❌ Failed to upload $(basename "$file"): $response${NC}"
        fi
        echo ""
    fi
done

rm -rf ./temp_docs

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}✅ Initialization complete!${NC}"
echo -e "${GREEN}📥 Downloaded: $downloaded/$total files${NC}"
echo -e "${GREEN}📤 Uploaded: $uploaded documents${NC}"
echo -e "${GREEN}========================================${NC}"