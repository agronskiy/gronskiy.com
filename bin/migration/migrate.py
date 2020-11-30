import os
import shutil

from glob import glob

TO_PATH = "/Users/agronsky/github/hugo-agronskiy/agronskiy/content/posts"

for lang in ["en", "ru"]:
    for path in glob(f"/Users/agronsky/github/agronskiy/_posts/{lang}/*.md"):
        filename = os.path.split(path)[-1].split(".")[0]
        new_dir = TO_PATH + f"/{lang}/{filename}"
        os.makedirs(new_dir, exist_ok=True)

        # print(f"Copy {path}, {new_dir + '/index.md'}")
        shutil.copyfile(path, new_dir + "/index.md")

        assets_dir = f"/Users/agronsky/github/agronskiy/assets/img/posts/{filename}"
        for asset in glob(f"{assets_dir}/*"):
            # print(f"Copy {asset}, {new_dir}")
            shutil.copy(asset, new_dir)
