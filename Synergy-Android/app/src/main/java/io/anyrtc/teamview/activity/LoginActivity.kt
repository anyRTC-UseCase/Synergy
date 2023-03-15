package io.anyrtc.teamview.activity

import android.content.Intent
import android.graphics.Color
import android.os.Bundle
import android.widget.Toast
import androidx.activity.viewModels
import androidx.lifecycle.lifecycleScope
import com.hjq.permissions.Permission
import com.hjq.permissions.XXPermissions
import io.anyrtc.teamview.databinding.ActivityLoginBinding
import io.anyrtc.teamview.utils.MyConstants
import io.anyrtc.teamview.utils.SpUtil
import io.anyrtc.teamview.vm.MainVM
import kotlinx.coroutines.launch

class LoginActivity : BaseActivity() {

    private lateinit var binding: ActivityLoginBinding
    private val vm: MainVM by viewModels()
    private var nickname = ""

    override fun onAttachedToWindow() {
        super.onAttachedToWindow()
        val uid = SpUtil.get().getString(MyConstants.UID, "")!!
        if (uid.isNotBlank()) {
            finish()
            startActivity(Intent(this, StartingActivity::class.java).also { it.putExtra("direct", true) })
        }
    }

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivityLoginBinding.inflate(layoutInflater)
        setContentView(binding.root)

        vm.loginResponse.observe(this) { result ->
            binding.loading.hideLoading()
            if (result.data == null) {
                Toast.makeText(applicationContext, result.errMsg, Toast.LENGTH_LONG).show()
                return@observe
            }

            val data = result.data
            startActivity(Intent(this, StartingActivity::class.java).also {
                it.putExtra("userTs", data.userTs)
                it.putExtra("userName", data.userName)
                it.putExtra("userType", data.userType)
                it.putExtra("workName", data.workName)
                it.putExtra(MyConstants.UID, data.uid)
            })
            finish()
        }
        binding.login.setOnClickListener {
            nickname = binding.nickname.text.toString()
            val workName = binding.work.text.toString()
            if (nickname.isBlank() || workName.isBlank()) {
                Toast.makeText(this, if (nickname.isBlank()) "请输入昵称" else "请输入工种", Toast.LENGTH_SHORT).show()
                return@setOnClickListener
            }

            showLoading()
            lifecycleScope.launch { vm.login(nickname, 1, workName) }
        }

        XXPermissions.with(this).permission(
            Permission.CAMERA,
            Permission.RECORD_AUDIO,
            Permission.WRITE_EXTERNAL_STORAGE,
            Permission.BLUETOOTH_CONNECT
        ).request { _, all ->
            if (!all) {
                Toast.makeText(this@LoginActivity, "请开启权限", Toast.LENGTH_SHORT).show()
                finish()
            }
        }
    }

    private fun showLoading() {
        binding.loading.run {
            setCardColor(Color.TRANSPARENT)
            setFontColor(Color.WHITE)
            showLoading("登录中")
        }
    }
}
