import JsBarcode from 'jsbarcode'

// Generate a Code128 barcode as a PNG data URL.
// Used for deposit pickup vouchers (scannable by pros) and admin deposit panels.
export function barcodeDataUrl(value, opts = {}) {
  if (!value) return ''
  const canvas = document.createElement('canvas')
  JsBarcode(canvas, String(value), {
    format: 'CODE128',
    displayValue: true,
    fontSize: 14,
    height: 60,
    margin: 6,
    background: '#ffffff',
    lineColor: '#001d32',
    ...opts,
  })
  return canvas.toDataURL('image/png')
}
