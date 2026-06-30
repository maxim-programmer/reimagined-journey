#!/bin/bash

echo "🔍 Checking PDF links..."
echo "================================"

urls=(
    "https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf"
    "https://www.learningcontainer.com/wp-content/uploads/2019/09/sample-pdf-file-for-testing.pdf"
    "https://www.clickdimensions.com/links/TestPDFfile.pdf"
    "https://www.orimi.com/pdf-test.pdf"
    "https://www.pdf995.com/samples/pdf.pdf"
    "https://www.africau.edu/images/default/sample.pdf"
    "https://www.antennahouse.com/antenna1/ACHOR001.pdf"
    "https://www.colorado.edu/engineering/sites/default/files/attached-files/example.pdf"
    "https://www.regulations.gov/docs/Training_CR_Web_Interface.pdf"
    "https://www.irs.gov/pub/irs-pdf/p1544.pdf"
)

for i in "${!urls[@]}"; do
    url="${urls[$i]}"
    echo -n "$((i+1)): "
    
    # Проверяем статус
    status=$(curl -o /dev/null -s -w "%{http_code}" -L "$url")
    
    # Проверяем размер
    size=$(curl -s -I -L "$url" | grep -i "content-length" | awk '{print $2}' | tr -d '\r')
    
    if [ "$status" = "200" ] || [ "$status" = "301" ] || [ "$status" = "302" ]; then
        if [ -n "$size" ] && [ "$size" -gt 1000 ]; then
            echo "✅ WORKING (Status: $status, Size: $((size/1024))KB)"
        else
            echo "⚠️  MAYBE WORKING (Status: $status, Size: unknown)"
        fi
    else
        echo "❌ FAILED (Status: $status)"
    fi
done

echo "================================"
