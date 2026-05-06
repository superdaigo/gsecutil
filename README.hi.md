# gsecutil - Google Secret Manager उपयोगिता

Google Secret Manager के लिए एक सरलीकृत कमांड-लाइन रैपर जो प्रति-परियोजना पासवर्ड मैनेजर की तरह काम करता है। सहज कमांड, क्लिपबोर्ड एकीकरण, संस्करण नियंत्रण, टीम-अनुकूल कॉन्फ़िगरेशन फ़ाइलों और ऑडिट लॉग के साथ सीक्रेट्स संग्रहीत, पुनर्प्राप्त और प्रबंधित करें।

## 🌍 भाषा संस्करण

- **English** - [README.md](README.md)
- **日本語** - [README.ja.md](README.ja.md)
- **中文** - [README.zh.md](README.zh.md)
- **Español** - [README.es.md](README.es.md)
- **हिंदी** - [README.hi.md](README.hi.md)（वर्तमान）
- **Português** - [README.pt.md](README.pt.md)

> **नोट**: सभी गैर-अंग्रेजी संस्करण मशीन-अनुवादित हैं। सबसे सटीक जानकारी के लिए, अंग्रेजी संस्करण देखें।

## त्वरित प्रारंभ

### स्थापना

अपने प्लेटफ़ॉर्म के लिए [रिलीज़ पेज](https://github.com/superdaigo/gsecutil/releases) से नवीनतम बाइनरी डाउनलोड करें, या Go के साथ स्थापित करें:

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

प्रत्येक प्रोजेक्ट में आमतौर पर अपनी कॉन्फ़िगरेशन फ़ाइल होती है जो प्रोजेक्ट ID, सीक्रेट नामकरण सम्मेलन और मेटाडेटा विशेषताओं को संग्रहीत करती है।

### 1. कॉन्फ़िगरेशन फ़ाइल बनाएं

एक कॉन्फ़िगरेशन फ़ाइल बनाने के लिए इंटरैक्टिव सेटअप चलाएं। यह आपसे आपका Google Cloud प्रोजेक्ट ID, सीक्रेट नाम उपसर्ग, डिफ़ॉल्ट सूची विशेषताएं और वैकल्पिक उदाहरण क्रेडेंशियल्स पूछेगा। जेनरेट की गई फ़ाइल डिफ़ॉल्ट रूप से वर्तमान निर्देशिका में `gsecutil.conf` के रूप में सहेजी जाती है (`--home` का उपयोग करके `~/.config/gsecutil/gsecutil.conf` में सहेजें)।

```bash
gsecutil config init
```

कॉन्फ़िगरेशन फ़ाइल इस क्रम में खोजी जाती है:
1. `--config` फ्लैग (यदि निर्दिष्ट है)
2. वर्तमान निर्देशिका: `gsecutil.conf`
3. होम निर्देशिका: `~/.config/gsecutil/gsecutil.conf`

### 2. सीक्रेट्स प्रबंधित करें

```bash
# एक सीक्रेट बनाएं
gsecutil create database-password

# नवीनतम संस्करण प्राप्त करें
gsecutil get database-password

# क्लिपबोर्ड में कॉपी करें
gsecutil get database-password --clipboard

# सभी सीक्रेट्स सूचीबद्ध करें
gsecutil list

# एक सीक्रेट अपडेट करें
gsecutil update database-password

# एक सीक्रेट हटाएं
gsecutil delete database-password
```

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

# क्रेडेंशियल मेटाडेटा (नाम सादे हैं — उपसर्ग स्वचालित रूप से जोड़ा जाता है)
credentials:
  - name: "database-password"    # "team-shared-database-password" तक पहुंचता है
    title: "Production Database Password"
    environment: "production"
    owner: "backend-team"
```

> **उपसर्ग पारदर्शी है:** जब उपसर्ग कॉन्फ़िगर किया जाता है, तो कमांड्स, कॉन्फ़िग और CSV फ़ाइलों में हमेशा सादे नामों का उपयोग करें। उपसर्ग स्वचालित रूप से जोड़ा और हटा दिया जाता है।

विस्तृत कॉन्फ़िगरेशन विकल्पों के लिए, [docs/configuration.md](docs/configuration.md) देखें।

## दस्तावेज़ीकरण

- **[कॉन्फ़िगरेशन गाइड](docs/configuration.md)** - विस्तृत कॉन्फ़िगरेशन विकल्प और उदाहरण
- **[कमांड संदर्भ](docs/commands.md)** - संपूर्ण कमांड दस्तावेज़ीकरण
- **[ऑडिट लॉगिंग सेटअप](docs/audit-logging.md)** - ऑडिट लॉग सक्षम और उपयोग करें
- **[समस्या निवारण गाइड](docs/troubleshooting.md)** - सामान्य समस्याएं और समाधान
- **[बिल्ड निर्देश](BUILD.md)** - स्रोत से बिल्ड करें
- **[विकास गाइड](WARP.md)** - WARP AI के साथ विकास

## लाइसेंस

यह प्रोजेक्ट MIT लाइसेंस के तहत लाइसेंस प्राप्त है - विवरण के लिए LICENSE फ़ाइल देखें।

## संबंधित

- [Google Cloud SDK](https://cloud.google.com/sdk)
- [Secret Manager दस्तावेज़ीकरण](https://cloud.google.com/secret-manager/docs)
