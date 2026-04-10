# Frontend Forms — Index

← [Back to Main README](../../../README.md) | [Frontend Overview →](../README.md) | [Conventions →](../conventions.md)

This directory contains implementation specs for every ERP screen: required backend endpoints, Valkey keys, shadcn/ui components, Zod schemas, and API request/response shapes.

---

## Form List

| # | Form | Phase | Description |
|---|------|-------|-------------|
| 01 | [Auth — Login](./01-auth-login.md) | 1 | Email OTP login, Google OAuth |
| 01b | [Auth — Email OTP Verify](./01b-email-otp-verify.md) | 1 | 6-digit code confirmation |
| 02 | [Employee Profile](./02-employee-profile.md) | 1 | Profile, GDP training, access flags |
| 03 | [Inbound Receiving](./03-inbound-receiving.md) | 1 | Goods receipt, quarantine, Inbound form |
| 04 | [Warehouse Zoning](./04-warehouse-zoning.md) | 1 | Zone management and storage rules |
| 05 | [Environment Log](./05-environment-log.md) | 3 | Temperature/humidity shift logging |
| 06 | [Assembly & Shipment (FEFO)](./06-assembly-shipment-fefo.md) | 2 | FEFO picking, MOS control, TTN |
| 07 | [Claims & Defects](./07-claims-defects.md) | 3 | Returns, recalls, Roszdravnadzor STOP signals |
| 08 | [Product Card](./08-product-card.md) | 1 | Master data for medicaments |
| 09 | [Inventory](./09-inventory.md) | 2 | Blind stocktaking |

---

## Related Documentation

| Document | Description |
|----------|-------------|
| [API Contract Index](../../api/README.md) | Full endpoint contract |
| [Database Schema](../../Diplom/migrate.md) | Full SQL schema |
| [Valkey Cache](../../valkey-cache.md) | Valkey data: OTP, cache, rate limiting |
| [Frontend Conventions](../conventions.md) | Tech stack and coding standards |

---

## 1. Common UI Elements

### 1.1 Header
- **Elements:** Logo, Navigation, Profile (UKEP status), Notification center.
- **Global search:** Search by INN, SKU, Batch number, **DataMatrix (KIZ)**, Barcode, Registration certificate #.

### 1.2 Footer
- **Legal data:** License number, expiry, name of the Authorized Person on shift.
- **Technical data:** Server status, DB version, GIS sync indicator (Chestny Znak / MDLP).

---

## 2. Access & Security

### 2.1 Auth — Login → [01-auth-login.md](./01-auth-login.md)
- **Login:** Email OTP / Google OAuth + mandatory **UKEP** binding for responsible persons.
- **Roles:** Admin, QP (Authorized Person), Warehouse Manager, Storekeeper, Pharmacist.
- **Access flags (admin):** Separate flag for **NS/PV** (Narcotic and Psychotropic substances) access.

### 2.1b Auth — Email OTP Verify → [01b-email-otp-verify.md](./01b-email-otp-verify.md)
- **Flow:** User enters 6-digit code from email. Code stored in Valkey (TTL 10 min, max 3 attempts).

### 2.2 Employee Profile → [02-employee-profile.md](./02-employee-profile.md)
- **Data:** Medical book scan, GDP training history, special zone access status.

---

## 3. Warehouse: Receiving & Quarantine

### 3.1 Inbound Receiving → [03-inbound-receiving.md](./03-inbound-receiving.md)
- **Main fields:** Supplier, Purchase type, Invoice #, Country, Manufacturer, **Registration certificate #**.
- **Batch data:** Batch #, Manufacture date, **Expiry date**.
- **Finance:** VAT rate, JNVLP markup control.
- **Quarantine status:** Received goods are system-blocked from all movements.
- **Release (admin/QP):** Acceptance protocol. Goods move to "available" only after QP signs off.

### 3.2 Warehouse Zoning → [04-warehouse-zoning.md](./04-warehouse-zoning.md)
- **Zones:** General, Cold chain (refrigerators), Flammable/alcohol, Safe zone (NS/PV, PKU).

### 3.3 Environment Log → [05-environment-log.md](./05-environment-log.md)
- **Log:** Daily entry form (2x/day): temperature and humidity per zone.

---

## 4. Logistics & Distribution

### 4.1 Assembly & Shipment (FEFO) → [06-assembly-shipment-fefo.md](./06-assembly-shipment-fefo.md)
- **FEFO logic:** System proposes batches with the **shortest remaining shelf life** (First Expired, First Out).
- **MOS control:** Shipment blocked if remaining shelf life < 60% (configurable).
- **Documentation:** Auto-generated quality certificate registry per TTN.
- **Cito!:** Urgent orders jump the assembly queue.

### 4.2 Claims & Defects → [07-claims-defects.md](./07-claims-defects.md)
- **Roszdravnadzor integration:** Auto-import of recalled batches. Matching batch → immediate STOP signal.
- **Returns:** Returns from pharmacies and to supplier, with photo evidence.

---

## 5. Master Data

### 5.1 Product Card → [08-product-card.md](./08-product-card.md)
- **Master data:** Photo, INN, **ATC classification**, Package multiplicity, dimensions/weight.
- **Flags:** JNVLP, MDLP (labeling), NS/PV (PKU), Cold chain.
- **Storage:** Temperature range and **max humidity**.

### 5.2 Inventory → [09-inventory.md](./09-inventory.md)
- **Mechanism:** Blind inventory (expected qty hidden until session closed).
- **Recount act (admin):** Netting surpluses and shortfalls within one price group.
- **Samples:** Tracking opened packages (control samples for the lab).

---

## 6. Analytics & Audit

### 6.1 Reporting
- **Expiry risk report:** 3/6/12-month financial risk analysis.
- **Turnover analysis:** ABC/XYZ analysis.

### 6.2 Audit (admin)
- **Versioning:** Full history of all quality certificates.
- **Immutable Logs:** SHA-256 chain hashing for every package action.

---

## Roadmap (Priority 2026)

### Phase 1 — "Survive" (Month 1-2)
- Auth, Users (Roles/Access flags).
- Product card (INN, registration #, NS/PV flags).
- Inbound → quarantine + expiry date.
- **Goal:** Start accounting without violating basic regulations.

### Phase 2 — "Regulatory Shield" (Month 3-4)
- **FEFO algorithm:** Replace standard FIFO.
- **DataMatrix labeling:** Integration with Chestny Znak.
- **JNVLP pricing control.**
- **Goal:** Protection from fines and batch expiry mix-ups.

### Phase 3 — "Safety & Quality" (Month 5-6)
- **UKEP (e-signature):** Acceptance and release protocols.
- **Environment logs:** Temperature and humidity.
- **Defect blocking:** Sync with recalled batch registries.
- **Goal:** Full compliance with GDP/GSP standards.

### Phase 4 — "Optimization" (Month 7+)
- **Address storage:** Racks/Cells.
- **Analytics:** ABC/XYZ and loss forecasting.
- **Logistics:** MOS (remaining shelf life) and dimension control.
- **Goal:** Maximize profit and reduce logistics costs.