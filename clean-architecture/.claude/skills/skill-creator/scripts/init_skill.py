#!/usr/bin/env python3
"""
スキル初期化ツール - テンプレートから新しいスキルを作成

使用法:
    init_skill.py <スキル名> --path <パス>

例:
    init_skill.py my-new-skill --path skills/public
    init_skill.py my-api-helper --path skills/private
    init_skill.py custom-skill --path /custom/location
"""

import sys
from pathlib import Path


SKILL_TEMPLATE = """---
name: {skill_name}
description: [TODO: スキルが何をするか、いつ使用するかについて完全で有益な説明を記述してください。このスキルを使用すべき具体的なシナリオ、ファイルタイプ、またはトリガーとなるタスクを含めてください。]
---

# {skill_title}

## 概要

[TODO: このスキルが何を可能にするかを1〜2文で説明してください]

## スキルの構造化

[TODO: このスキルの目的に最も適した構造を選択してください。一般的なパターン：

**1. ワークフローベース** (順次プロセスに最適)
- 明確なステップバイステップの手順がある場合に適している
- 例: DOCXスキルで「ワークフロー決定木」→「読み取り」→「作成」→「編集」
- 構造: ## 概要 → ## ワークフロー決定木 → ## ステップ1 → ## ステップ2...

**2. タスクベース** (ツールコレクションに最適)
- スキルが異なる操作/機能を提供する場合に適している
- 例: PDFスキルで「クイックスタート」→「PDFを結合」→「PDFを分割」→「テキストを抽出」
- 構造: ## 概要 → ## クイックスタート → ## タスクカテゴリ1 → ## タスクカテゴリ2...

**3. リファレンス/ガイドライン** (標準や仕様に最適)
- ブランドガイドライン、コーディング標準、要件に適している
- 例: ブランドスタイリングで「ブランドガイドライン」→「カラー」→「タイポグラフィ」→「機能」
- 構造: ## 概要 → ## ガイドライン → ## 仕様 → ## 使用方法...

**4. 機能ベース** (統合システムに最適)
- スキルが複数の相互関連する機能を提供する場合に適している
- 例: プロダクトマネジメントで「コア機能」→番号付き機能リスト
- 構造: ## 概要 → ## コア機能 → ### 1. 機能 → ### 2. 機能...

パターンは必要に応じて組み合わせることができます。ほとんどのスキルはパターンを組み合わせています（例：タスクベースで始め、複雑な操作にはワークフローを追加）。

完了したら、この「スキルの構造化」セクション全体を削除してください - これは単なるガイダンスです。]

## [TODO: 選択した構造に基づいて最初のメインセクションに置き換えてください]

[TODO: ここにコンテンツを追加してください。既存のスキルの例を参照：
- 技術スキルのコードサンプル
- 複雑なワークフローの決定木
- 現実的なユーザーリクエストを含む具体例
- 必要に応じてscripts/templates/referencesへの参照]

## リソース

このスキルには、異なるタイプのバンドルリソースの整理方法を示すリソースディレクトリの例が含まれています：

### scripts/
特定の操作を実行するために直接実行できる実行可能コード（Python/Bashなど）。

**他のスキルからの例:**
- PDFスキル: `fill_fillable_fields.py`, `extract_form_field_info.py` - PDF操作のユーティリティ
- DOCXスキル: `document.py`, `utilities.py` - ドキュメント処理用のPythonモジュール

**適している用途:** Pythonスクリプト、シェルスクリプト、または自動化、データ処理、特定の操作を実行する実行可能コード。

**注意:** スクリプトはコンテキストに読み込まずに実行される場合がありますが、パッチ適用や環境調整のためにClaudeによって読み取られる可能性があります。

### references/
Claudeのプロセスと思考を知らせるために、コンテキストに読み込まれることを意図したドキュメントと参照資料。

**他のスキルからの例:**
- プロダクトマネジメント: `communication.md`, `context_building.md` - 詳細なワークフローガイド
- BigQuery: APIリファレンスドキュメントとクエリ例
- Finance: スキーマドキュメント、企業ポリシー

**適している用途:** 詳細なドキュメント、APIリファレンス、データベーススキーマ、包括的なガイド、またはClaudeが作業中に参照すべき詳細情報。

### assets/
コンテキストに読み込まれることを意図していない、Claudeが生成する出力内で使用されるファイル。

**他のスキルからの例:**
- ブランドスタイリング: PowerPointテンプレートファイル（.pptx）、ロゴファイル
- フロントエンドビルダー: HTML/Reactボイラープレートプロジェクトディレクトリ
- タイポグラフィ: フォントファイル（.ttf, .woff2）

**適している用途:** テンプレート、ボイラープレートコード、ドキュメントテンプレート、画像、アイコン、フォント、または最終出力でコピーまたは使用されるファイル。

---

**不要なディレクトリは削除できます。** すべてのスキルが3種類のリソースすべてを必要とするわけではありません。
"""

EXAMPLE_SCRIPT = '''#!/usr/bin/env python3
"""
{skill_name}のサンプルヘルパースクリプト

これは直接実行できるプレースホルダースクリプトです。
実際の実装に置き換えるか、不要な場合は削除してください。

他のスキルの実際のスクリプト例:
- pdf/scripts/fill_fillable_fields.py - PDFフォームフィールドを埋める
- pdf/scripts/convert_pdf_to_images.py - PDFページを画像に変換
"""

def main():
    print("{skill_name}のサンプルスクリプトです")
    # TODO: 実際のスクリプトロジックをここに追加してください
    # データ処理、ファイル変換、API呼び出しなど

if __name__ == "__main__":
    main()
'''

EXAMPLE_REFERENCE = """# {skill_title}のリファレンスドキュメント

これは詳細なリファレンスドキュメントのプレースホルダーです。
実際のリファレンスコンテンツに置き換えるか、不要な場合は削除してください。

他のスキルの実際のリファレンスドキュメント例:
- product-management/references/communication.md - ステータス更新の包括的ガイド
- product-management/references/context_building.md - コンテキスト収集の詳細
- bigquery/references/ - APIリファレンスとクエリ例

## リファレンスドキュメントが有用な場合

リファレンスドキュメントは以下に最適です：
- 包括的なAPIドキュメント
- 詳細なワークフローガイド
- 複雑な複数ステップのプロセス
- メインのSKILL.mdには長すぎる情報
- 特定のユースケースでのみ必要なコンテンツ

## 構造の提案

### APIリファレンスの例
- 概要
- 認証
- エンドポイントと例
- エラーコード
- レート制限

### ワークフローガイドの例
- 前提条件
- ステップバイステップの指示
- 一般的なパターン
- トラブルシューティング
- ベストプラクティス
"""

EXAMPLE_ASSET = """# サンプルアセットファイル

このプレースホルダーは、アセットファイルが保存される場所を表しています。
実際のアセットファイル（テンプレート、画像、フォントなど）に置き換えるか、不要な場合は削除してください。

アセットファイルはコンテキストに読み込まれることを意図しておらず、
Claudeが生成する出力内で使用されます。

他のスキルのアセットファイル例:
- ブランドガイドライン: logo.png, slides_template.pptx
- フロントエンドビルダー: HTML/Reactボイラープレートのhello-world/ディレクトリ
- タイポグラフィ: custom-font.ttf, font-family.woff2
- データ: sample_data.csv, test_dataset.json

## 一般的なアセットタイプ

- テンプレート: .pptx, .docx, ボイラープレートディレクトリ
- 画像: .png, .jpg, .svg, .gif
- フォント: .ttf, .otf, .woff, .woff2
- ボイラープレートコード: プロジェクトディレクトリ、スターターファイル
- アイコン: .ico, .svg
- データファイル: .csv, .json, .xml, .yaml

注: これはテキストプレースホルダーです。実際のアセットは任意のファイルタイプが可能です。
"""


def title_case_skill_name(skill_name):
    """Convert hyphenated skill name to Title Case for display."""
    return ' '.join(word.capitalize() for word in skill_name.split('-'))


def init_skill(skill_name, path):
    """
    テンプレートから新しいスキルディレクトリを初期化します。

    Args:
        skill_name: スキルの名前
        path: スキルディレクトリを作成するパス

    Returns:
        作成されたスキルディレクトリのパス、またはエラーの場合はNone
    """
    # スキルディレクトリパスを決定
    skill_dir = Path(path).resolve() / skill_name

    # ディレクトリがすでに存在するかチェック
    if skill_dir.exists():
        print(f"❌ エラー: スキルディレクトリがすでに存在します: {skill_dir}")
        return None

    # スキルディレクトリを作成
    try:
        skill_dir.mkdir(parents=True, exist_ok=False)
        print(f"✅ スキルディレクトリを作成しました: {skill_dir}")
    except Exception as e:
        print(f"❌ ディレクトリ作成エラー: {e}")
        return None

    # テンプレートからSKILL.mdを作成
    skill_title = title_case_skill_name(skill_name)
    skill_content = SKILL_TEMPLATE.format(
        skill_name=skill_name,
        skill_title=skill_title
    )

    skill_md_path = skill_dir / 'SKILL.md'
    try:
        skill_md_path.write_text(skill_content)
        print("✅ SKILL.mdを作成しました")
    except Exception as e:
        print(f"❌ SKILL.md作成エラー: {e}")
        return None

    # サンプルファイル付きのリソースディレクトリを作成
    try:
        # サンプルスクリプト付きのscripts/ディレクトリを作成
        scripts_dir = skill_dir / 'scripts'
        scripts_dir.mkdir(exist_ok=True)
        example_script = scripts_dir / 'example.py'
        example_script.write_text(EXAMPLE_SCRIPT.format(skill_name=skill_name))
        example_script.chmod(0o755)
        print("✅ scripts/example.pyを作成しました")

        # サンプルリファレンスドキュメント付きのreferences/ディレクトリを作成
        references_dir = skill_dir / 'references'
        references_dir.mkdir(exist_ok=True)
        example_reference = references_dir / 'api_reference.md'
        example_reference.write_text(EXAMPLE_REFERENCE.format(skill_title=skill_title))
        print("✅ references/api_reference.mdを作成しました")

        # サンプルアセットプレースホルダー付きのassets/ディレクトリを作成
        assets_dir = skill_dir / 'assets'
        assets_dir.mkdir(exist_ok=True)
        example_asset = assets_dir / 'example_asset.txt'
        example_asset.write_text(EXAMPLE_ASSET)
        print("✅ assets/example_asset.txtを作成しました")
    except Exception as e:
        print(f"❌ リソースディレクトリ作成エラー: {e}")
        return None

    # 次のステップを表示
    print(f"\n✅ スキル '{skill_name}' を {skill_dir} に正常に初期化しました")
    print("\n次のステップ:")
    print("1. SKILL.mdを編集してTODO項目を完成させ、説明を更新してください")
    print("2. scripts/, references/, assets/のサンプルファイルをカスタマイズまたは削除してください")
    print("3. スキル構造を確認するために、準備ができたらバリデーターを実行してください")

    return skill_dir


def main():
    if len(sys.argv) < 4 or sys.argv[2] != '--path':
        print("使用法: init_skill.py <スキル名> --path <パス>")
        print("\nスキル名の要件:")
        print("  - ハイフン区切りの識別子（例: 'data-analyzer'）")
        print("  - 小文字、数字、ハイフンのみ使用可能")
        print("  - 最大40文字")
        print("  - ディレクトリ名と正確に一致する必要があります")
        print("\n例:")
        print("  init_skill.py my-new-skill --path skills/public")
        print("  init_skill.py my-api-helper --path skills/private")
        print("  init_skill.py custom-skill --path /custom/location")
        sys.exit(1)

    skill_name = sys.argv[1]
    path = sys.argv[3]

    print(f"🚀 スキルを初期化中: {skill_name}")
    print(f"   場所: {path}")
    print()

    result = init_skill(skill_name, path)

    if result:
        sys.exit(0)
    else:
        sys.exit(1)


if __name__ == "__main__":
    main()
