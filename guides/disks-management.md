# Datakeeper Node CLI Disk Management Guide

## Overview
The Disk Manager provides essential storage management capabilities for node operations, allowing Datakeepers to manage data storage paths, relocate volumes, and adjust storage limits for licensed nodes.

## Command: disks add

**Purpose:** Register a new data storage path for a licensed node.

**Flags:**
* `--license` (required): License identifier
* `--path` (required): Filesystem path for data storage
* `--limit` (required): Storage limit in gigabytes
* `--address` (required): Node address identifier

**Usage:**
```
./denode disks add --address 0x1...0 --license 12345 --path /var/lib/denode/storage --limit 10
```

## Command: disks move

**Purpose:** Relocate an existing storage volume to a new filesystem location.

**Flags:**
* `--license` (required): License identifier
* `--path` (required): Current storage path
* `--new-path` (required): Destination storage path
* `--address` (required): Node address identifier

**Usage:**
```
./denode disks move --address 0x1...0 --license 12345 --path /var/lib/denode/old-storage --new-path /var/lib/denode/new-storage
```

## Command: disks resize

**Purpose:** Modify the storage capacity limit for an existing volume.

**Flags:**
* `--license` (required): License identifier
* `--path` (required): Target storage path
* `--limit` (required): New storage limit in gigabytes
* `--address` (required): Node address identifier

**Usage:**
```
./denode disks resize --address 0x1...0 --license 12345 --path /var/lib/denode/storage --limit 20
```
