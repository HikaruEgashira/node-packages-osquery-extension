Output demo-unique.gif

Set FontSize 24
Set Width 1200
Set Height 600
Set Padding 20
Set TypingSpeed 30ms

Hide
Type "osqueryi --extension $(which node_packages_extension)" Enter
Sleep 1s
Show

Type "SELECT DISTINCT name, version FROM node_packages ORDER BY name;" Sleep 100ms Enter

Sleep 1.5s

Hide
Type ".quit" Enter
