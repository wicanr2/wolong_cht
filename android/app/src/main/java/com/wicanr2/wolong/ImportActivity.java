package com.wicanr2.wolong;

import android.app.Activity;
import android.content.Intent;
import android.content.res.AssetManager;
import android.database.Cursor;
import android.net.Uri;
import android.os.Bundle;
import android.provider.DocumentsContract;
import android.view.Gravity;
import android.view.View;
import android.widget.Button;
import android.widget.LinearLayout;
import android.widget.TextView;

import java.io.File;
import java.io.FileOutputStream;
import java.io.InputStream;
import java.io.OutputStream;

/**
 * 匯入原版資料。
 *
 * <p>⚠ Android 11 以上，使用者選的資料夾給的是 {@code content://} URI，
 * 而遊戲本體用的是 {@code os.ReadFile(路徑)}——**兩者接不起來**。
 * 所以要有這一步：把檔案複製進 app 的私有目錄，之後 Go 端就只是讀檔案。
 *
 * <p>不用 androidx 的 DocumentFile，改用 {@link DocumentsContract} 直接查，
 * 少一個相依就少一次「建置環境要能上網」的門檻。
 */
public final class ImportActivity extends Activity {

    /** 挑資料夾的請求碼。 */
    private static final int PICK_TREE = 1;

    /**
     * 字型檔的檔名。它們與遊戲資料**放在不同的目錄**：
     * 遊戲資料進 {@code orig/}、字型進 {@code eten/}。
     *
     * <p>⚠ 倚天字型不隨程式散布（deny-list 的邊界對手機版一樣成立），
     * 所以它也要使用者自備。同一次匯入裡認出來就分過去，
     * 使用者不必知道這兩個目錄的差別。
     */
    private static final String[] FONT_FILES = {
        "STDFONT.15", "ASCFONT.15", "SPCFONT.15",
        "STDFONT.24", "ASCFONT.24",
    };

    /** 缺了它就不是原版資料夾——用它當「挑對了沒」的判準。 */
    private static final String REQUIRED = "SINARIO.DAT";

    /**
     * APK 內嵌原版資料的 assets 根目錄（docs/spec/72 §4）。
     *
     * <p>底下是 {@code orig/} 與 {@code eten/} 兩個子目錄，與解開後
     * 在 {@code getFilesDir()} 裡的版面**一模一樣**——兩邊同名，
     * 之後改動只要改一處。
     *
     * <p>⚠ **不是每個建置都有它。** 內嵌由 {@code WOLONG_BUNDLE_DATA}
     * 在建置時決定，預設不內嵌；沒有內嵌時這個目錄不存在，
     * 流程要安靜地退回 SAF 匯入，不能當成錯誤。
     */
    private static final String BUNDLE_ROOT = "gamedata";

    /** 解開時的暫存後綴。**先解完再改名**，中斷不會留下半套資料。 */
    private static final String PART_SUFFIX = ".part";

    private TextView status;

    @Override
    protected void onCreate(Bundle state) {
        super.onCreate(state);
        if (hasData()) {
            startGame();
            return;
        }
        // ⭐ 內嵌資料的建置：解開就走，使用者不必挑資料夾。
        if (unpackBundled()) {
            startGame();
            return;
        }
        LinearLayout root = new LinearLayout(this);
        root.setOrientation(LinearLayout.VERTICAL);
        root.setGravity(Gravity.CENTER);
        root.setPadding(48, 48, 48, 48);

        TextView title = new TextView(this);
        title.setText("臥龍傳 Remake");
        title.setTextSize(28);
        title.setGravity(Gravity.CENTER);
        root.addView(title);

        status = new TextView(this);
        status.setText("請選擇原版資料夾（要有 " + REQUIRED + "）。\n"
            + "倚天點陣字放同一個資料夾也可以，會自動分開。\n"
            + "原版資料不隨程式散布，請自備。");
        status.setTextSize(16);
        status.setGravity(Gravity.CENTER);
        status.setPadding(0, 32, 0, 32);
        root.addView(status);

        Button pick = new Button(this);
        pick.setText("選擇資料夾");
        pick.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                startActivityForResult(
                    new Intent(Intent.ACTION_OPEN_DOCUMENT_TREE), PICK_TREE);
            }
        });
        root.addView(pick);
        setContentView(root);
    }

    @Override
    protected void onActivityResult(int request, int result, Intent data) {
        super.onActivityResult(request, result, data);
        if (request != PICK_TREE || result != RESULT_OK || data == null
                || data.getData() == null) {
            return;
        }
        try {
            int n = copyTree(data.getData());
            if (!new File(getFilesDir(), "orig/" + REQUIRED).exists()) {
                // ⚠ 缺檔要**指名**：「匯入失敗」對使用者毫無用處。
                status.setText("這個資料夾裡沒有 " + REQUIRED + "（複製了 " + n
                    + " 個檔）。請選原版遊戲的資料夾。");
                return;
            }
            startGame();
        } catch (Exception e) {
            status.setText("匯入失敗：" + e);
        }
    }

    /**
     * 從 APK 內嵌的 assets 解開原版資料，回傳資料是不是因此就齊了。
     *
     * <p>沒有內嵌（一般建置）就回 {@code false}，由 SAF 那條路接手。
     *
     * <p>⚠ **先解到 {@code .part} 再改名**。半途被殺掉的話，
     * 留下的是一個沒有人會去讀的 {@code orig.part}，而不是一個
     * 「{@code SINARIO.DAT} 在、其他檔缺一半」的 {@code orig}——
     * 後者會讓 {@link #hasData()} 回 true，之後每次開機都直接進遊戲然後爆掉。
     */
    private boolean unpackBundled() {
        AssetManager am = getAssets();
        try {
            String[] top = am.list(BUNDLE_ROOT);
            if (top == null || top.length == 0) {
                return false;
            }
            for (String sub : new String[] {"orig", "eten"}) {
                String[] names = am.list(BUNDLE_ROOT + "/" + sub);
                if (names == null || names.length == 0) {
                    continue;
                }
                File part = new File(getFilesDir(), sub + PART_SUFFIX);
                deleteTree(part);
                if (!part.mkdirs()) {
                    return false;
                }
                for (String name : names) {
                    copyAsset(am, BUNDLE_ROOT + "/" + sub + "/" + name,
                        new File(part, name));
                }
                File dst = new File(getFilesDir(), sub);
                deleteTree(dst);
                if (!part.renameTo(dst)) {
                    return false;
                }
            }
            return hasData();
        } catch (Exception e) {
            return false;
        }
    }

    /** 遞迴刪一棵目錄樹。只用在 app 自己的私有目錄底下。 */
    private static void deleteTree(File f) {
        File[] kids = f.listFiles();
        if (kids != null) {
            for (File k : kids) {
                deleteTree(k);
            }
        }
        f.delete();
    }

    private static void copyAsset(AssetManager am, String name, File dst)
            throws Exception {
        InputStream in = am.open(name);
        try {
            OutputStream out = new FileOutputStream(dst);
            try {
                byte[] buf = new byte[64 * 1024];
                for (int k = in.read(buf); k > 0; k = in.read(buf)) {
                    out.write(buf, 0, k);
                }
            } finally {
                out.close();
            }
        } finally {
            in.close();
        }
    }

    /** 資料齊不齊：`orig/SINARIO.DAT` 在就算齊。 */
    private boolean hasData() {
        return new File(getFilesDir(), "orig/" + REQUIRED).exists();
    }

    private void startGame() {
        startActivity(new Intent(this, MainActivity.class));
        finish();
    }

    /** 把選中資料夾裡的檔案複製進 app 私有目錄，回傳複製了幾個。 */
    private int copyTree(Uri tree) throws Exception {
        File orig = new File(getFilesDir(), "orig");
        File eten = new File(getFilesDir(), "eten");
        orig.mkdirs();
        eten.mkdirs();

        Uri children = DocumentsContract.buildChildDocumentsUriUsingTree(
            tree, DocumentsContract.getTreeDocumentId(tree));
        Cursor c = getContentResolver().query(children, new String[] {
            DocumentsContract.Document.COLUMN_DOCUMENT_ID,
            DocumentsContract.Document.COLUMN_DISPLAY_NAME,
            DocumentsContract.Document.COLUMN_MIME_TYPE,
        }, null, null, null);
        if (c == null) {
            throw new IllegalStateException("讀不到資料夾內容");
        }
        int n = 0;
        try {
            while (c.moveToNext()) {
                String id = c.getString(0);
                String name = c.getString(1);
                String mime = c.getString(2);
                if (DocumentsContract.Document.MIME_TYPE_DIR.equals(mime)) {
                    continue; // 只收最上層的檔案，子目錄不遞迴
                }
                File dst = new File(isFont(name) ? eten : orig, name);
                copy(DocumentsContract.buildDocumentUriUsingTree(tree, id), dst);
                n++;
            }
        } finally {
            c.close();
        }
        return n;
    }

    private static boolean isFont(String name) {
        for (String f : FONT_FILES) {
            if (f.equalsIgnoreCase(name)) {
                return true;
            }
        }
        return false;
    }

    private void copy(Uri src, File dst) throws Exception {
        InputStream in = getContentResolver().openInputStream(src);
        if (in == null) {
            throw new IllegalStateException("開不了 " + src);
        }
        try {
            OutputStream out = new FileOutputStream(dst);
            try {
                byte[] buf = new byte[64 * 1024];
                for (int k = in.read(buf); k > 0; k = in.read(buf)) {
                    out.write(buf, 0, k);
                }
            } finally {
                out.close();
            }
        } finally {
            in.close();
        }
    }
}
