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
      <h1>Backup 3DS dashboard</h1>
    </header>

    <div class="count"><span></span>Backups of <% $dir_name %></div>
    <div class="count">Used space <span><%% du -hs  "$dir" | awk '{print $1}'   %></span> </div>

    <div class="list">

      %
        for backup in $(find "$dir" -name '*.zip' ); do
          backup_name="$(basename "$backup")"
          echo "<div class='item'>
                  <div class='item-icon' style='background:rgba(200,255,94,0.12)'>💾</div>
                  <div class='item-body'>
                    <div class='item-name'>$backup_name</div>
                    <div class='item-desc'>$(stat -c '%y' "$backup" )</div>
                  </div>
                </div>"
        done
      %

    </div>
  </div>
</body>
</html>
