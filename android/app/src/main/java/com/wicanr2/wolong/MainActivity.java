package com.wicanr2.wolong;

import android.app.Activity;
import android.os.Bundle;

import java.io.File;

import com.wicanr2.wolong.mobile.wolongmobile.EbitenView;
import com.wicanr2.wolong.mobile.wolongmobile.Wolongmobile;

import go.Seq;

public final class MainActivity extends Activity {
    @Override
    protected void onCreate(Bundle state) {
        super.onCreate(state);
        Seq.setContext(getApplicationContext());
        Wolongmobile.setDataRoot(dataRoot());
        // 驗收參數走 Intent extra——**app 不繼承 adb 的環境變數**。
        String frames = getIntent().getStringExtra("fp_frames");
        if (frames != null) {
            Wolongmobile.setFingerprintFrames(frames);
        }
        setContentView(new EbitenView(this));
    }

    /**
     * 返回鍵：先交給遊戲關掉開著的面板，沒東西可關才交回系統。
     *
     * <p>⚠ 即時制沒有暫停，返回鍵直接結束的話誤觸一次就是進度沒了。
     */
    @Override
    public void onBackPressed() {
        if (Wolongmobile.back()) {
            return;
        }
        super.onBackPressed();
    }

    /**
     * 原版資料與點陣字的根目錄：{@code <root>/orig} 放 69 個原版檔、
     * {@code <root>/eten} 放字型。兩者都由使用者自備，不隨 APK 散布。
     *
     * <p>優先用 app 內部目錄（{@code getFilesDir}）：匯入流程把檔案複製到
     * 那裡，而它的擁有權單純是 app 自己，讀得到是確定的。
     *
     * <p>⚠ 外部目錄（{@code getExternalFilesDir}）只當退路。它看起來很方便
     * —— {@code adb push} 進得去 —— 但 Android 11 以上那條路徑是 FUSE 掛的，
     * 由 adb 寫進去的檔案 app **讀不到**（{@code permission denied}），
     * 而目錄本身還是列得出來，所以症狀像「檔案在那裡但壞掉」。
     */
    private String dataRoot() {
        File internal = getFilesDir();
        if (new File(internal, "orig").isDirectory()) {
            return internal.getAbsolutePath();
        }
        File ext = getExternalFilesDir(null);
        if (ext != null) {
            return ext.getAbsolutePath();
        }
        return internal.getAbsolutePath();
    }
}
