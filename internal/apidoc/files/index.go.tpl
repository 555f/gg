<!doctype html>
<html lang="en" class="dark">
  <head>
    <meta charset="UTF-8" />
    <link rel="icon" href="/favicon.ico" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title></title>

    <script>window.apiData = {{ .APIData }};</script>    
    <style rel="stylesheet" crossorigin>{{ .Style }}</style>
  </head>
  <body class="bg-white antialiased dark:bg-zinc-900">
    <div id="app"></div>
  </body>
  <script type="module" crossorigin>{{ .Script }}</script>
</html>
