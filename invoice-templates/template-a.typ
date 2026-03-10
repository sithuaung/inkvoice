#let invoice-number = sys.inputs.at("invoice-number", default: "INK-0000")
#let issue-date = sys.inputs.at("issue-date", default: "")
#let due-date = sys.inputs.at("due-date", default: "")
#let client-name = sys.inputs.at("client-name", default: "")
#let client-email = sys.inputs.at("client-email", default: "")
#let client-company = sys.inputs.at("client-company", default: "")
#let company-name = sys.inputs.at("company-name", default: "Inkvoice")
#let company-email = sys.inputs.at("company-email", default: "")
#let subtotal = sys.inputs.at("subtotal", default: "0")
#let tax-total = sys.inputs.at("tax-total", default: "0")
#let total = sys.inputs.at("total", default: "0")
#let currency = sys.inputs.at("currency", default: "USD")
#let notes = sys.inputs.at("notes", default: "")
#let items-json = sys.inputs.at("items", default: "[]")

#set page(margin: 2cm)
#set text(font: "Inter", size: 10pt, fill: luma(30))

#let fmt-money(cents-str) = {
  let cents = int(cents-str)
  let dollars = calc.floor(cents / 100)
  let remainder = calc.rem(calc.abs(cents), 100)
  [#currency #str(dollars).#if remainder < 10 [0]#str(remainder)]
}

// Header
#grid(
  columns: (1fr, 1fr),
  align: (left, right),
  [
    #text(size: 24pt, weight: "bold", fill: rgb("#2563eb"))[INVOICE]
    #v(4pt)
    #text(size: 12pt)[#("\#") #invoice-number]
  ],
  [
    #text(weight: "bold")[#company-name]
    #if company-email != "" [\ #company-email]
  ],
)

#v(1.5cm)

// Dates and client info
#grid(
  columns: (1fr, 1fr),
  [
    #text(weight: "bold")[Bill To:]
    #v(4pt)
    #client-name
    #if client-company != "" [\ #client-company]
    #if client-email != "" [\ #client-email]
  ],
  align(right)[
    #text(weight: "bold")[Issue Date:] #issue-date \
    #text(weight: "bold")[Due Date:] #due-date
  ],
)

#v(1cm)

// Items table
#let items = json.decode(items-json)

#table(
  columns: (1fr, auto, auto, auto),
  inset: 8pt,
  stroke: 0.5pt + luma(200),
  fill: (_, y) => if y == 0 { rgb("#2563eb").lighten(90%) },
  table.header(
    text(weight: "bold")[Description],
    text(weight: "bold")[Qty],
    text(weight: "bold")[Price],
    text(weight: "bold")[Amount],
  ),
  ..for item in items {
    (
      item.at("description", default: ""),
      str(item.at("quantity", default: 0)),
      fmt-money(str(item.at("unit_price", default: 0))),
      fmt-money(str(item.at("amount", default: 0))),
    )
  }
)

#v(0.5cm)

// Totals
#align(right)[
  #grid(
    columns: (auto, auto),
    column-gutter: 1cm,
    row-gutter: 6pt,
    [Subtotal:], fmt-money(subtotal),
    [Tax:], fmt-money(tax-total),
    text(weight: "bold", size: 12pt)[Total:], text(weight: "bold", size: 12pt)[#fmt-money(total)],
  )
]

// Notes
#if notes != "" [
  #v(1cm)
  #line(length: 100%, stroke: 0.5pt + luma(200))
  #v(0.5cm)
  #text(weight: "bold")[Notes:]
  #v(4pt)
  #notes
]
