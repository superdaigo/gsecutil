# gsecutil - Google Secret Manager उपयोगिता

🚀 कॉन्फ़िगरेशन फ़ाइल समर्थन और टीम-अनुकूल सुविधाओं के साथ Google Secret Manager के लिए एक सरलीकृत कमांड-लाइन रैपर।

## 🌍 भाषा संस्करण

- **English** - [README.md](README.md)
- **日本語** - [README.ja.md](README.ja.md)
- **中文** - [README.zh.md](README.zh.md)
- **Español** - [README.es.md](README.es.md)
- **हिंदी** - [README.hi.md](README.hi.md) (वर्तमान)
- **Português** - [README.pt.md](README.pt.md)

> **नोट**: सभी गैर-अंग्रेजी संस्करण मशीन-अनुवादित हैं। सबसे सटीक जानकारी के लिए, अंग्रेजी संस्करण देखें।

## त्वरित प्रारंभ

### स्थापना

[रिलीज़ पेज](https://github.com/superdaigo/gsecutil/releases) से अपने प्लेटफ़ॉर्म के लिए नवीनतम बाइनरी डाउनलोड करें:

```bash
# macOS Apple Silicon
curl -L https://github.com/superdaigo/gsecutil/releases/latest/download/gsecutil-darwin-arm64 -o gsecutil
chmod +x gsecutil
sudo mv gsecutil /usr/local/bin/

# macOS Intel
curl -L https://github.com/superdaigo/gsecutil/releases/latest/download/gsecutil-darwin-amd64 -o gsecutil
chmod +x gsecutil
sudo mv gsecutil /usr/local/bin/

# Linux
curl -L https://github.com/superdaigo/gsecutil/releases/latest/download/gsecutil-linux-amd64 -o gsecutil
chmod +x gsecutil
sudo mv gsecutil /usr/local/bin/

# Windows (PowerShell)
Invoke-WebRequest -Uri "https://github.com/superdaigo/gsecutil/releases/latest/download/gsecutil-windows-amd64.exe" -OutFile "gsecutil.exe"
# Move to a directory in your PATH, e.g., C:\Windows\System32
Move-Item gsecutil.exe C:\Windows\System32\gsecutil.exe
```

या Go के साथ स्थापित करें:
```bash
go install github.com/superdaigo/gsecutil@latest
```

### आवश्यक शर्तें

- [Google Cloud SDK (gcloud)](https://cloud.google.com/sdk/docs/install) स्थापित और प्रमाणित
- Secret Manager API सक्षम के साथ Google Cloud प्रोजेक्ट

### प्रमाणीकरण

```bash
# gcloud के साथ प्रमाणित करें
gcloud auth login

# डिफ़ॉल्ट प्रोजेक्ट सेट करें
gcloud config set project YOUR_PROJECT_ID

# या पर्यावरण चर सेट करें
export GSECUTIL_PROJECT=YOUR_PROJECT_ID
```

## मूल उपयोग

### एक सीक्रेट बनाएं
```bash
# इंटरैक्टिव इनपुट
gsecutil create database-password

# कमांड लाइन से
gsecutil create api-key -d "sk-1234567890"

# फ़ाइल से
gsecutil create config --data-file ./config.json
```

### एक सीक्रेट प्राप्त करें
```bash
# नवीनतम संस्करण प्राप्त करें
gsecutil get database-password

# क्लिपबोर्ड में कॉपी करें
gsecutil get api-key --clipboard

# विशिष्ट संस्करण प्राप्त करें
gsecutil get api-key --version 2
```

### सीक्रेट्स सूचीबद्ध करें
```bash
# सभी सीक्रेट्स सूचीबद्ध करें
gsecutil list

# लेबल द्वारा फ़िल्टर करें
gsecutil list --filter "labels.env=prod"
```

### एक सीक्रेट अपडेट करें
```bash
# इंटरैक्टिव इनपुट
gsecutil update database-password

# कमांड लाइन से
gsecutil update api-key -d "new-secret-value"
```

### एक सीक्रेट हटाएं
```bash
gsecutil delete old-secret
```

## कॉन्फ़िगरेशन

gsecutil प्रोजेक्ट-विशिष्ट सेटिंग्स के लिए कॉन्फ़िगरेशन फ़ाइलों का समर्थन करता है। कॉन्फ़िगरेशन फ़ाइलें इस क्रम में खोजी जाती हैं:

1. `--config` फ्लैग (यदि निर्दिष्ट है)
2. वर्तमान डाइरेक्टरी: `gsecutil.conf`
3. होम डाइरेक्टरी: `~/.config/gsecutil/gsecutil.conf`

### कॉन्फ़िगरेशन उदाहरण

```yaml
# प्रोजेक्ट ID (यदि पर्यावरण या gcloud के माध्यम से सेट है तो वैकल्पिक)
project: "my-project-id"

# टीम संगठन के लिए सीक्रेट नाम उपसर्ग
prefix: "team-shared-"

# list कमांड में प्रदर्शित करने के लिए डिफ़ॉल्ट विशेषताएं
list:
  attributes:
    - title
    - owner
    - environment

# क्रेडेंशियल मेटाडेटा (नाम सादे हैं — प्रीफ़िक्स स्वचालित रूप से जोड़ा जाता है)
credentials:
  - name: "database-password"    # "team-shared-database-password" तक पहुंचता है
    title: "Production Database Password"
    environment: "production"
    owner: "backend-team"
```

> **प्रीफ़िक्स पारदर्शी है:** जब प्रीफ़िक्स कॉन्फ़िगर किया जाता है, तो कमांड, कॉन्फ़िग और CSV फ़ाइलों में हमेशा सादे नाम (बिना प्रीफ़िक्स के) का उपयोग करें। प्रीफ़िक्स स्वचालित रूप से जोड़ा और हटाया जाता है।

### त्वरित प्रारंभ

```bash
# इंटरैक्टिव रूप से कॉन्फ़िगरेशन उत्पन्न करें
gsecutil config init

# या प्रोजेक्ट-विशिष्ट कॉन्फ़िग बनाएं
echo 'project: "my-project-123"' > gsecutil.conf
```

विस्तृत कॉन्फ़िगरेशन विकल्पों के लिए, [docs/configuration.md](docs/configuration.md) देखें।

## मुख्य विशेषताएं

- ✅ **सरल CRUD संचालन** - सीक्रेट्स प्रबंधित करने के लिए सहज कमांड
- ✅ **क्लिपबोर्ड एकीकरण** - सीक्रेट्स को सीधे क्लिपबोर्ड में कॉपी करें
- ✅ **संस्करण प्रबंधन** - विशिष्ट संस्करणों तक पहुंच और संस्करण जीवनचक्र प्रबंधित करें
- ✅ **कॉन्फ़िगरेशन फ़ाइल समर्थन** - टीम-अनुकूल मेटाडेटा और संगठन
- ✅ **एक्सेस प्रबंधन** - बुनियादी IAM नीति प्रबंधन
- ✅ **ऑडिट लॉग** - देखें कि किसने कब सीक्रेट्स तक पहुंच बनाई
- ✅ **एकाधिक इनपुट विधियां** - इंटरैक्टिव, इनलाइन, या फ़ाइल-आधारित
- ✅ **क्रॉस-प्लेटफ़ॉर्म** - Linux, macOS, Windows (amd64/arm64)

## दस्तावेज़ीकरण

- **[कॉन्फ़िगरेशन गाइड](docs/configuration.md)** - विस्तृत कॉन्फ़िगरेशन विकल्प और उदाहरण
- **[कमांड संदर्भ](docs/commands.md)** - संपूर्ण कमांड दस्तावेज़ीकरण
- **[ऑडिट लॉगिंग सेटअप](docs/audit-logging.md)** - ऑडिट लॉग सक्षम और उपयोग करें
- **[समस्या निवारण गाइड](docs/troubleshooting.md)** - सामान्य समस्याएं और समाधान
- **[बिल्ड निर्देश](BUILD.md)** - स्रोत से बिल्ड करें
- **[विकास गाइड](WARP.md)** - WARP AI के साथ विकास

## सामान्य कमांड

```bash
# सीक्रेट विवरण दिखाएं
gsecutil describe my-secret

# संस्करण इतिहास दिखाएं
gsecutil describe my-secret --show-versions

# ऑडिट लॉग देखें
gsecutil auditlog my-secret

# एक्सेस प्रबंधित करें
gsecutil access list my-secret
gsecutil access grant my-secret --principal user:alice@example.com

# कॉन्फ़िगरेशन सत्यापित करें
gsecutil config validate

# कॉन्फ़िगरेशन दिखाएं
gsecutil config show
```

## लाइसेंस

यह प्रोजेक्ट MIT लाइसेंस के तहत लाइसेंस प्राप्त है - विवरण के लिए LICENSE फ़ाइल देखें।

## संबंधित

- [Google Cloud SDK](https://cloud.google.com/sdk)
- [Secret Manager दस्तावेज़ीकरण](https://cloud.google.com/secret-manager/docs)
