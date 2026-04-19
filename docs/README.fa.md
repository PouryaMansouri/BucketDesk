# BucketDesk

**BucketDesk** یک اپ سبک و دو زبانه برای مدیریت bucketهای MinIO و S3-compatible است؛ بدون این‌که لازم باشد MinIO Console را در اختیار کاربران قرار بدهید.

این ابزار برای شرکت‌هایی مناسب است که می‌خواهند کاربران بتوانند فایل آپلود کنند، bucket را مرور کنند، لینک objectها را کپی کنند و objectهای مجاز را حذف کنند، اما دسترسی مدیریتی کامل MinIO را نداشته باشند.

English: [README](../README.md)

## قابلیت‌ها

- رابط کاربری دو زبانه: فارسی و انگلیسی.
- تغییر خودکار جهت صفحه بین RTL و LTR.
- چند پروفایل اتصال MinIO/S3.
- تنظیم Endpoint، Bucket، Region، CDN URL و Path-style.
- تست اتصال و بررسی دسترسی نوشتن روی bucket.
- مرور bucket مثل پوشه‌ها با S3 prefix.
- آپلود چند فایل در مسیر فعلی.
- انتخاب و حذف objectها.
- کپی URL عمومی objectها.
- اجرای محلی، بدون سرویس خارجی، دیتابیس یا telemetry.

## چرا BucketDesk؟

در بسیاری از شرکت‌ها، MinIO Console را به کاربران عادی نمی‌دهند، چون قدرت عملیاتی زیادی دارد. BucketDesk یک پنل محدود و تمیز برای مدیریت objectها می‌دهد و ادمین می‌تواند دسترسی واقعی را با policyهای S3 محدود کند.

## مدل امنیتی پیشنهادی

از credential ریشه MinIO استفاده نکنید.

برای هر تیم یا جریان کاری یک access key جدا بسازید و دسترسی آن را فقط به bucket و prefixهای لازم محدود کنید. نمونه policyها در این فایل آمده‌اند:

[IAM policy examples](./IAM_POLICIES.md)

## تکنولوژی

- **Go** برای بک‌اند و ساخت خروجی سبک و قابل نصب روی سیستم‌عامل‌های مختلف.
- **React + TypeScript** برای UI سریع و دو زبانه.
- **AWS SDK for Go v2** برای APIهای سازگار با S3.

Go انتخاب شده چون می‌شود برنامه را به شکل یک فایل اجرایی برای macOS، Windows و Linux منتشر کرد و کاربر نهایی لازم نیست Node.js نصب کند.

## توسعه

نیازمندی‌ها:

- Go 1.23+
- Node.js 20+
- npm

نصب dependencyها:

```bash
npm install
go mod download
```

اجرای بک‌اند:

```bash
go run ./cmd/bucketdesk
```

اجرای UI در حالت توسعه:

```bash
npm run dev:web
```

در حالت توسعه، Vite درخواست‌های `/api` را به سرور Go proxy می‌کند.

## ساخت خروجی

```bash
npm run build:app
```

خروجی در این مسیر ساخته می‌شود:

```text
dist/bucketdesk
```

اجرا:

```bash
./dist/bucketdesk
```

برنامه یک آدرس محلی چاپ می‌کند، معمولاً:

```text
http://127.0.0.1:5217
```

## بسته‌های قابل انتشار

BucketDesk می‌تواند در سه حالت ساده منتشر شود:

| سیستم‌عامل | نصب‌کننده | پرتابل |
| --- | --- | --- |
| Windows | فایل setup با پسوند `.exe` | فایل `.zip` شامل `bucketdesk.exe` |
| macOS | فایل `.dmg` شامل `BucketDesk.app` | فایل `.tar.gz` شامل باینری |
| Linux | فایل `.deb` | فایل `.tar.gz` شامل باینری |

حالت پرتابل نیاز به نصب ندارد. کاربر فایل را extract می‌کند و برنامه را اجرا می‌کند. برنامه به صورت خودکار سرور محلی را بالا می‌آورد و مرورگر را باز می‌کند.

ساخت آرشیوهای پرتابل:

```bash
VERSION=v0.1.0 ./scripts/package-portable.sh
```

ساخت DMG برای مک:

```bash
VERSION=v0.1.0 ARCH=arm64 ./scripts/package-macos-dmg.sh
VERSION=v0.1.0 ARCH=amd64 ./scripts/package-macos-dmg.sh
```

ساخت `.deb` روی لینوکس:

```bash
VERSION=v0.1.0 ARCH=amd64 ./scripts/package-linux-deb.sh
VERSION=v0.1.0 ARCH=arm64 ./scripts/package-linux-deb.sh
```

نصب‌کننده ویندوز در GitHub Actions با Inno Setup ساخته می‌شود.

## انتشار نسخه

برای ساخت Release عمومی در GitHub:

```bash
git tag v0.1.0
git push origin v0.1.0
```

Workflow انتشار این فایل‌ها را می‌سازد:

- نصب‌کننده `.exe` برای Windows
- فایل `.dmg` برای مک Intel و Apple Silicon
- فایل `.deb` برای Linux amd64 و arm64
- نسخه‌های پرتابل برای Windows، macOS و Linux

## مجوز

BucketDesk تحت Apache License 2.0 منتشر می‌شود. فایل [LICENSE](../LICENSE) را ببینید.
