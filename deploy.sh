systemctl stop onthego
git pull
npx @tailwindcss/cli -i static/input.css -o static/styles.css
systemctl start onthego