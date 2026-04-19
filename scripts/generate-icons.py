#!/usr/bin/env python3
from pathlib import Path
from PIL import Image, ImageDraw

ROOT = Path(__file__).resolve().parent.parent
ASSETS = ROOT / "assets"
ASSETS.mkdir(exist_ok=True)

SVG = """<svg width="512" height="512" viewBox="0 0 512 512" fill="none" xmlns="http://www.w3.org/2000/svg">
  <rect width="512" height="512" rx="112" fill="#0F172A"/>
  <path d="M132 196C132 163.967 157.967 138 190 138H322C354.033 138 380 163.967 380 196V316C380 348.033 354.033 374 322 374H190C157.967 374 132 348.033 132 316V196Z" fill="#197BFF"/>
  <path d="M168 210C168 188.461 185.461 171 207 171H305C326.539 171 344 188.461 344 210V303C344 324.539 326.539 342 305 342H207C185.461 342 168 324.539 168 303V210Z" fill="#FFFFFF" fill-opacity="0.95"/>
  <path d="M204 236C204 220.536 216.536 208 232 208H280C295.464 208 308 220.536 308 236V282C308 297.464 295.464 310 280 310H232C216.536 310 204 297.464 204 282V236Z" fill="#197BFF"/>
  <path d="M230 251H282" stroke="#FFFFFF" stroke-width="18" stroke-linecap="round"/>
  <path d="M256 225V293" stroke="#FFFFFF" stroke-width="18" stroke-linecap="round"/>
  <path d="M185 138L211 102H301L327 138H185Z" fill="#00A663"/>
</svg>
"""


def draw_icon(size: int) -> Image.Image:
    scale = size / 512
    img = Image.new("RGBA", (size, size), (0, 0, 0, 0))
    d = ImageDraw.Draw(img)

    def box(values):
        return tuple(round(v * scale) for v in values)

    def radius(value):
        return round(value * scale)

    d.rounded_rectangle(box((0, 0, 512, 512)), radius=radius(112), fill="#0F172A")
    d.rounded_rectangle(box((132, 138, 380, 374)), radius=radius(58), fill="#197BFF")
    d.rounded_rectangle(box((168, 171, 344, 342)), radius=radius(39), fill=(255, 255, 255, 242))
    d.rounded_rectangle(box((204, 208, 308, 310)), radius=radius(28), fill="#197BFF")
    d.rounded_rectangle(box((185, 102, 327, 138)), radius=radius(16), fill="#00A663")
    d.line(box((230, 251, 282, 251)), fill="#FFFFFF", width=max(1, radius(18)))
    d.line(box((256, 225, 256, 293)), fill="#FFFFFF", width=max(1, radius(18)))
    return img


def main():
    (ASSETS / "bucketdesk.svg").write_text(SVG, encoding="utf-8")
    png512 = draw_icon(512)
    png512.save(ASSETS / "bucketdesk.png")

    ico_sizes = [16, 24, 32, 48, 64, 128, 256]
    png512.save(ASSETS / "bucketdesk.ico", sizes=[(size, size) for size in ico_sizes])

    iconset = ASSETS / "BucketDesk.iconset"
    if iconset.exists():
        for child in iconset.iterdir():
            child.unlink()
    else:
        iconset.mkdir()

    icon_sizes = [
        (16, "icon_16x16.png"),
        (32, "icon_16x16@2x.png"),
        (32, "icon_32x32.png"),
        (64, "icon_32x32@2x.png"),
        (128, "icon_128x128.png"),
        (256, "icon_128x128@2x.png"),
        (256, "icon_256x256.png"),
        (512, "icon_256x256@2x.png"),
        (512, "icon_512x512.png"),
        (1024, "icon_512x512@2x.png"),
    ]
    for size, name in icon_sizes:
        draw_icon(size).save(iconset / name)


if __name__ == "__main__":
    main()
