package com.wicanr2.wolong;

import android.app.Activity;
import android.os.Bundle;

import com.wicanr2.wolong.mobile.wolongmobile.EbitenView;

import go.Seq;

public final class MainActivity extends Activity {
    @Override
    protected void onCreate(Bundle state) {
        super.onCreate(state);
        Seq.setContext(getApplicationContext());
        setContentView(new EbitenView(this));
    }
}
