package io.anyrtc.teamview

import android.app.Application
import android.util.Log
import com.hjq.permissions.XXPermissions
import com.kongzue.dialog.util.DialogSettings
import io.anyrtc.teamview.utils.MyConstants
import io.anyrtc.teamview.utils.SpUtil
import okhttp3.OkHttpClient
import rxhttp.RxHttpPlugins
import rxhttp.wrapper.ssl.HttpsUtils
import java.util.concurrent.TimeUnit
import javax.net.ssl.SSLSession
import kotlin.properties.Delegates

class App : Application() {

    companion object{
        var app : App by Delegates.notNull()
        var token: String = ""
    }

    override fun onCreate() {
        super.onCreate()
        app = this
        SpUtil.init(this)
        DialogSettings.style = DialogSettings.STYLE.STYLE_IOS
        RxHttpPlugins.init(getDefaultOkHttpClient()).setDebug(BuildConfig.DEBUG)
            .setOnParamAssembly {
                val url = it.url
                Log.e("::", "setOnParamAssembly:token=$token, url=$url")
                if (token.isNotEmpty()) {
                    it.addHeader("Authorization", "Bearer $token")
                    //it.addHeader("Artc-Token", token)
                }
                it
            }
    }

    private fun getDefaultOkHttpClient(): OkHttpClient {
        val sslParams = HttpsUtils.getSslSocketFactory()
        return OkHttpClient.Builder()
            .connectTimeout(6, TimeUnit.SECONDS)
            .readTimeout(6, TimeUnit.SECONDS)
            .writeTimeout(6, TimeUnit.SECONDS)
            .sslSocketFactory(sslParams.sSLSocketFactory, sslParams.trustManager) //添加信任证书
            .hostnameVerifier { _: String?, _: SSLSession? -> true } //忽略host验证
            .build()
    }
}
