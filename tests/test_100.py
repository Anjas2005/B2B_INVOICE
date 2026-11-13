import json
import random
import uuid
from datetime import date, timedelta
from urllib.error import HTTPError, URLError
from urllib.request import Request, urlopen

URL = "https://b2b-invoice.onrender.com/invoices"

customers = [
    "ABC Traders",
    "XYZ Stores",
    "Global Enterprises",
    "Reliance Supplies",
    "Tech Solutions",
    "Prime Industries",
    "Apex Corporation",
    "Sunrise Retail",
    "Metro Distributors",
    "Nova Technologies",
]

batch_id = uuid.uuid4().hex[:8]
invoices = []

for i in range(1, 101):
    required = random.choice(
        [1000, 2500, 5000, 10000, 15000, 25000, 50000, 75000, 100000]
    )

    paid = random.choice([0, required // 4, required // 2, required])

    due_date = date.today() + timedelta(days=random.randint(1, 180))

    invoices.append(
        {
            "customer": random.choice(customers),
            "invoice_no": f"BULK-{batch_id}-{i:04d}",
            "required_amount": required,
            "paid_amount": paid,
            "due_date": due_date.isoformat(),
        }
    )

# Serialize all 1000 invoices into one JSON array
payload = json.dumps(invoices).encode("utf-8")

# Save a copy for inspection
with open("invoices_1000.json", "wb") as f:
    f.write(payload)

print(f"Generated: {len(invoices)} invoices")
print(f"Payload size: {len(payload) / 1024:.2f} KB")
print("Sending request...")

request = Request(
    URL, data=payload, headers={"Content-Type": "application/json"}, method="POST"
)

try:
    with urlopen(request, timeout=120) as response:
        print("HTTP Status:", response.status)
        print("Response:", response.read().decode())

except HTTPError as e:
    print("HTTP Error:", e.code)
    print(e.read().decode())

except URLError as e:
    print("Connection Error:", e.reason)
