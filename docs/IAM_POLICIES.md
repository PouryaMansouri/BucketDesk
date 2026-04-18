# IAM / MinIO Policy Examples

These examples show how to give BucketDesk users limited bucket access without exposing the MinIO Console.

## Read and Write a Single Prefix

Replace:

- `media-production` with your bucket name.
- `uploads/team-a/` with the allowed prefix.

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": ["s3:ListBucket"],
      "Resource": ["arn:aws:s3:::media-production"],
      "Condition": {
        "StringLike": {
          "s3:prefix": ["uploads/team-a/*", "uploads/team-a/"]
        }
      }
    },
    {
      "Effect": "Allow",
      "Action": ["s3:GetObject", "s3:PutObject", "s3:DeleteObject"],
      "Resource": ["arn:aws:s3:::media-production/uploads/team-a/*"]
    }
  ]
}
```

## Read-Only Prefix

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": ["s3:ListBucket"],
      "Resource": ["arn:aws:s3:::media-production"],
      "Condition": {
        "StringLike": {
          "s3:prefix": ["public/*", "public/"]
        }
      }
    },
    {
      "Effect": "Allow",
      "Action": ["s3:GetObject"],
      "Resource": ["arn:aws:s3:::media-production/public/*"]
    }
  ]
}
```

## فارسی

این policyها برای این هستند که کاربر بتواند فقط روی bucket یا prefix مشخص‌شده کار کند و نیازی به MinIO Console نداشته باشد.

برای استفاده امن‌تر:

- credential ریشه MinIO را وارد BucketDesk نکنید.
- برای هر تیم یک access key جدا بسازید.
- دسترسی `s3:DeleteObject` را فقط وقتی بدهید که واقعاً لازم است.
- اگر کاربر فقط باید فایل‌ها را ببیند، policy read-only کافی است.
