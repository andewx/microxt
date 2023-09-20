folder1 = new Folder(1, "Microcontrollers")
folder2 = new Folder(2, "Radar Hardware")
folder3 = new Folder(3, "Profiles")
folder4 = new Folder(4, "Layouts")
folder5 = new Folder(5, "Analysis")
folder6 = new Folder(6, "Design Tools")
file1 = new File(1, "Radar Beam K-LD7 24GHZ")
file2 = new File(2, "Texas Instruments 60GHZ")
file3 = new File(3, "Infineon 24GHZ Transceiver")
file4 = new File(4, "Velocity Profile")
file5 = new File(5, "Range Profile")
file6 = new File(6, "Doppler Profile")
file7 = new File(7, "Range-Doppler Heatmap")
file8 = new File(8, "Range-Angle Heatmap")
file9 = new File(9, "Range-Doppler-Angle Heatmap")
file10 = new File(10, "Configure Radar Experiment")
file11 = new File(11, "Max Speed/Distance Analysis")
file12 = new File(12, "Chirp Design Tool")
file13 = new File(13, "Antenna Design Tool")
file14 = new File(14, "FSK Design Tool")
file15 = new File(15, "Arduino Portenta H7")
file16 = new File(16, "Arduino MKR Wifi 1010")
file17 = new File(17, "Arduino Nano 33 BLE")


folder1.children.push(file15)
folder1.children.push(file16)
folder1.children.push(file17)

folder2.children.push(file1)
folder2.children.push(file2)
folder2.children.push(file3)

folder3.children.push(file4)
folder3.children.push(file5)
folder3.children.push(file6)

folder4.children.push(file7)
folder4.children.push(file8)
folder4.children.push(file9)


folder5.children.push(file10)
folder5.children.push(file11)

folder6.children.push(file12)
folder6.children.push(file13)
folder6.children.push(file14)


fileview = new FileView()
fileview.children.push(folder1)
fileview.children.push(folder2)
fileview.children.push(folder3)
fileview.children.push(folder4)
fileview.children.push(folder5)
fileview.children.push(folder6)

divider = new Divider(0, "divider")
mApp.add(divider)
mApp.add(fileview)

mApp.init()

mApp.render()
