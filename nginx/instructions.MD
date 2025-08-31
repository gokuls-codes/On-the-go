1. Install nginx
2. `sudo systemctl start nginx`
3. `sudo systemctl enable nginx`
4. Navigate to `/etc/nginx/sites-available/`
5. Write the contents of `onthego.conf` to `onthego.conf` there
6. Create a symlink from available to enabled using `sudo ln -s /etc/nginx/sites-available/onthego.conf /etc/nginx/sites-enabled/`
7. Test config using `sudo nginx -t`
8. Reload nginx using `sudo systemctl reload nginx`
9.
