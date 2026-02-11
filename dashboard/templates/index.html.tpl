<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>Backup 3DS dashboard</title>
  <link href="https://fonts.googleapis.com/css2?family=Syne:wght@400;600;700;800&family=DM+Mono:wght@300;400;500&display=swap" rel="stylesheet" />
  <link rel="stylesheet" href="style.css">
</head>
<body>
  <div class="glow"></div>
  <div class="container">

    <header>
      <a href="/"><h1>Backup 3DS dashboard</h1></a>
    </header>

    <div class="count"><span><%% ls "$BACKUP_DEST" | wc -l  %></span> Consoles</div>
    <div class="count">Used space <span><%% du -hs  "$BACKUP_DEST" | awk '{print $1}'   %></span> </div>

    <div class="list">

      %
        for console in $(find "$BACKUP_DEST/"  -maxdepth 1 -type d -not -path "$BACKUP_DEST/" ); do
          console_name="$(basename "$console")"
          echo "<a href='$console_name.html'>
                  <div class='item'>
                    <div class='item-icon' style='background:rgba(200,255,94,0.12)'>💾</div>
                    <div class='item-body'>
                      <div class='item-name'>$console_name</div>
                    </div>
                  </div>
                </a>"
        done
      %

    </div>
  </div>
</body>
</html>
