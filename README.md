# Baggsy Project

### Description
```
Here’s a concise description of the Baggsy app based on the functionality:

Baggsy: A Bag Registration and Tracking App

Overview

Baggsy is a Go (Golang) and React-based application designed to register, track, and manage agricultural bags. It utilizes QR codes to identify each bag uniquely and supports both Parent (main) bags and Child (sub) bags. Baggsy also handles linking parent bags to bill IDs—ensuring that each parent bag is used exactly once—and allows scanning of child bag QR codes to look up associated bills.

Key Features
	1.	Parent Bag Registration
	•	A parent bag is registered by scanning or entering its QR code.
	•	The QR code includes a child count that indicates how many child bags can be associated with this parent.
	•	The parent bag’s details (e.g., QRCode, BagType, ChildCount) are stored in the bags table.
	2.	Child Bag Registration
	•	After a parent bag is created, you can register child bags one by one.
	•	Each child bag also has a unique QR code.
	•	On the backend, each child bag references its parentBag field, effectively linking it to the parent bag.
	•	A child count limit ensures the number of child bags does not exceed the capacity set by the parent bag’s QR code.
	3.	Linking Parent Bags to Bills
	•	A parent bag can be linked to a Bill ID, marking it so it cannot be reused.
	•	The backend sets the parent bag as linked or soft-deletes it (depending on the configuration) once it’s tied to a bill.
	4.	Scanning Child Bags for Bill Lookup
	•	A child bag can be scanned to quickly retrieve the Bill ID that its parent bag is associated with.
	•	Ideal for tracing which bill a child bag belongs to.
	5.	Automatic Window Closing & UI Flow
	•	Once all child bags are registered for a parent (based on the child count), the app automatically resets to the parent registration window—offering a smooth user flow.

Architecture & Technology Stack
	1.	Frontend:
	•	React (JavaScript) with functional components and hooks.
	•	Axios for HTTP requests.
	•	A simple, clean UI that provides two main forms:
	•	Parent Bag Registration: For creating a parent bag and specifying capacity (childCount).
	•	Child Bag Registration: For scanning child bags and linking them to a parent.
	2.	Backend:
	•	Go (Golang) using Gin as the web framework.
	•	GORM (Go ORM) for database operations.
	•	PostgreSQL (or another relational DB) to store bag and linking data.
	•	Endpoints:
	•	/register-bag: Registers parent or child bags.
	•	/link-child-bag: Links a scanned child bag to its parent bag (if separate from /register-bag).
	•	/link-bag-to-bill: Ties a parent bag to a bill.
	3.	Database Schema:
	•	bags Table: Stores QR codes, bag types (parent or child), child count, and an optional parentBag field for child bags.
	<!-- •	bag_maps Table (optional approach): Maps parent bag to child bag, if the system needs a separate table for relationships. -->
	•	links Table (optional approach): Links a parent bag to a bill ID.

Typical Use Cases
	1.	Field Use / Warehouse:
	•	Worker scans a parent bag QR code, the system registers how many child bags it can hold.
	•	Worker registers child bags (scanning each child’s QR code), automatically tying them to the parent bag.
	2.	Billing / Tracking:
	•	Once a parent bag is ready to be billed, the bag is linked to a Bill ID.
	•	Scanning a child bag in the future allows quick retrieval of the parent’s Bill ID.

Benefits & Purpose
	•	Traceability: Ensure each parent bag’s content (child bags) is well tracked.
	•	Data Integrity: Each parent bag can only be used once (linked to a single bill).
	•	Scalability: Simple flow for registering, linking, and searching.
	•	User Experience: Automatically resets to parent registration after child bag capacities are met, streamlining usage for on-the-ground workers.

Baggsy thus provides an end-to-end solution for QR-based agricultural bag registration and traceability, ensuring smooth scanning, bill linkage, and limited child bag capacities per parent bag.
```

#### Backend
Run the backend:
```bash
go run main.go
```

#### Frontend
Run the frontend:
```bash
npm start
```

#### Docker
Start the project using Docker:
```bash
docker-compose up --build
```

