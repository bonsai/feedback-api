import os
import requests

API = "https://api.github.com/repos/bonsai/meikyoku-janken/issues"


def create_feedback(day: int, theme: str, message: str, category: str = "デザイン"):
    token = os.environ["GITHUB_TOKEN"]
    body = f"## デザインフィードバック\n\n- DAY: {day}\n- お題: {theme}\n- カテゴリ: {category}\n\n### 内容\n\n{message.strip()}"
    r = requests.post(API, headers={"Authorization": f"Bearer {token}", "Accept": "application/vnd.github+json", "X-GitHub-Api-Version": "2022-11-28"}, json={"title": f"[Design Feedback] DAY {day} {theme}", "body": body, "labels": ["design"]}, timeout=10)
    r.raise_for_status()
    return r.json()["html_url"]


if __name__ == "__main__":
    print(create_feedback(3, "色", "お題をもっと大きくしてほしい"))
