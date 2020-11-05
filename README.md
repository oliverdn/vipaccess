# PSD2 born, accessibility suffers
Electron UI of Symantec VIP Access app to get one time passwords.

![before after](https://i.imgur.com/qQNtIKV.png)

## Command line usage:

Print help. ```Symatec.exe -h```

Generate an id and secret pair, requires internet connection. ```Symantec.exe -issuer="Symantec"```

Create OTP from secret, it works offline. ```Symantec.exe -show=EJM7URT4OI4DXCVN6G5TQ7KGPM5ABCDE```

## UI:

The program downloads electron files to AppData\Roaming directory.

There is a generated ```key.txt``` file next to Symantec.exe executable, create a backup,  
visit www.paypal.com/cgi-bin/webscr?cmd=_setup-security-key, register your new id.
