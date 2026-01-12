# DeNode Manager Desktop Installation and Launch Guide

## Step 1. Download DeNode Manager Desktop

1. Visit:  
   [https://github.com/DeNetPRO/Node/releases](https://github.com/DeNetPRO/Node/releases)
2. Download the latest application executable for your OS (Linux or Windows).

![Assets](/assets/asset_list_1.png)

---

## Step 2. Install and start the app

### Windows

#### 1. Open the downloaded **.exe** file.
#### 2. Follow the installer instructions:  
   **Next → Install → Finish**

![Screen](/assets/windows_screen_1.png)

#### 3. Launch the **DeNode Manager** application.

![Screen](/assets/windows_screen_2.png)

---

### Linux

#### 1. Download the latest version of DeNode Manager Desktop

**Option 1 — via GitHub (see Step 1).**

**Option 2 — via Terminal:**
Open Terminal: Use Ctrl + Alt + T or your terminal shortcut. Or SSH to your remote host.

Then run the command:

```bash
curl -LO https://github.com/DeNetPRO/Node/releases/download/dm-desktop-v1.0.0/DeNode_Manager_1.0.0_amd64.deb
```

---

#### 2. Install the application

Install the downloaded `.deb` package:

```bash
sudo dpkg -i DeNode_Manager_1.0.0_amd64.deb
```

If dependency errors occur, run:

```bash
sudo apt-get install -f -y
```

This command will automatically fix missing dependencies and complete the installation.

---

#### 3. Allow **denode** executable

```bash
sudo chmod +x /usr/lib/DeNode\ Manager/bin/denode
```

> ⚠ The path may vary depending on your system.

---

![Screen](/assets/linux_screen_1.jpg)

#### 4. Start the app

Launch the application — open your applications menu and find DeNode Manager.

If your system provides a launcher shortcut, you can use it as well.

![Screen](/assets/linux_screen_2.jpg)
