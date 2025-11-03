Output demo.gif

Set FontSize 24
Set Width 1200
Set Height 600
Set Padding 20
Set TypingSpeed 30ms

Type "osqueryi --extension $(which node_packages_extension)" Sleep 200ms Enter

Sleep 1.5s

Type "SELECT name, version, manager FROM node_packages LIMIT 5;" Sleep 200ms Enter

Sleep 2s

Type "SELECT manager, COUNT(*) as count FROM node_packages GROUP BY manager;" Sleep 200ms Enter

Sleep 2s

Type "SELECT name, version FROM node_packages WHERE name LIKE '%react%' LIMIT 5;" Sleep 200ms Enter

Sleep 2s

Type ".quit" Sleep 200ms Enter

Sleep 500ms
