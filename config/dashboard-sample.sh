cat > wazi-samples.json <<EOF
[
  {
    "displayName": "IBM Developer for z/OS on VS Code",
    "description": "Development stack for z/OS application development for COBOL, PL/I, HLASM, REXX, JCL",
    "tags": ["cobol", "pl1", "hlasm", "rexx", "jcl", "zOS", "ibm"],
    "url": "https://github.com/IBM/zopeneditor-sample/tree/devfile",
    "icon": {
      "base64data": "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",
      "mediatype": "image/png"
    }
  }
]
EOF