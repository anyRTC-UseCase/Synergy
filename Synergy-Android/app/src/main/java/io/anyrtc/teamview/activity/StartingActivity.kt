package io.anyrtc.teamview.activity

import android.content.Intent
import android.graphics.Color
import android.os.Bundle
import android.view.View
import android.widget.Toast
import androidx.activity.viewModels
import com.hjq.permissions.Permission
import com.hjq.permissions.XXPermissions
import io.anyrtc.teamview.databinding.ActivityStartingBinding
import io.anyrtc.teamview.utils.MyConstants
import io.anyrtc.teamview.utils.RtcManager
import io.anyrtc.teamview.utils.SpUtil
import io.anyrtc.teamview.vm.MainVM
import java.util.*

class StartingActivity : BaseActivity() {

    private lateinit var binding: ActivityStartingBinding
    private val vm: MainVM by viewModels()

    private val onlineStatusCallback by lazy {
        object : RtcManager.RtcHandler() {
            override fun onOtherClientLogin() {
                runOnUiThread {
                    SpUtil.edit { it.remove(MyConstants.UID) }
                    vm.removeRtcHandler(this)
                    Toast.makeText(this@StartingActivity, "已在其他手机登录", Toast.LENGTH_SHORT).show()
                    startActivity(Intent(this@StartingActivity, LoginActivity::class.java)/*.also {
					    it.addFlags(Intent.FLAG_ACTIVITY_CLEAR_TOP)
                    }*/)
                    finish()
                }
            }
        }
    }

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivityStartingBinding.inflate(layoutInflater)
        setContentView(binding.root)

        val isDirect = intent.getBooleanExtra("direct", false)
        val bundle = intent.extras
        val newBundle = Bundle()

        val spGet = SpUtil.get()
        val userName = spGet.getString(MyConstants.USER_NAME, "")!!
        val workName = spGet.getString(MyConstants.WORK_NAME, "")!!
        var loginSuccess = false
        var userClick = false

        if (isDirect || bundle == null) {
            showLoading("登录中")

            vm.loginResponse.observe(this) { response ->
                if (response.data == null) {
                    /*binding.root.postDelayed({
                        vm.login(userName, 1, workName)
                    }, 200)*/
                    binding.loading.hideLoading()
                    Toast.makeText(this, response.errMsg, Toast.LENGTH_LONG).show()
                    return@observe
                }

                loginSuccess = true
                val data = response.data
                newBundle.run {
                    putInt("userTs", data.userTs)
                    putString("userName", data.userName)
                    putInt("userType", data.userType)
                    putString("workName", data.workName)
                    putString(MyConstants.UID, data.uid)
                }
                if (userClick) {
                    showLoading("加载中")
                    vm.createRoom(true)
                    userClick = false
                    return@observe
                }
                registerOnlineStatus()
                binding.loading.hideLoading()
            }
            vm.login(userName, 1, workName)
        } else {
            registerOnlineStatus()
            loginSuccess = true
            newBundle.putAll(bundle)
        }

        vm.createRoomResponse.observe(this) {
            if (it.data == null) {
                binding.loading.hideLoading()
                Toast.makeText(this, it.errMsg/*"create room failed: ${it.errMsg}"*/, Toast.LENGTH_SHORT).show()
                return@observe
            }

            val newIntent = Intent(this, MainActivity::class.java)
            newIntent.putExtras(newBundle)
            val data = it.data
            newIntent.run {
                putExtra("rtcToken", data.rtcToken)
                putExtra("roomId", data.roomId)
                putExtra("roomName", data.roomName)
                putExtra("roomTs", data.roomTs)
            }

            startActivity(newIntent)
            binding.loading.hideLoading()
        }

        binding.start.setOnClickListener {
            if (binding.loading.visibility == View.VISIBLE) {
                return@setOnClickListener
            }

            if (!loginSuccess) {
                showLoading("登录中")
                userClick = true
                vm.login(userName, 1, workName)
                return@setOnClickListener
            }
            showLoading("加载中")
            vm.createRoom(true)
        }

        XXPermissions.with(this).permission(
            Permission.CAMERA,
            Permission.RECORD_AUDIO,
            Permission.WRITE_EXTERNAL_STORAGE
        ).request { _, all ->
            if (!all) {
                Toast.makeText(this@StartingActivity, "请开启权限", Toast.LENGTH_SHORT).show()
                finish()
            }
        }
    }

    override fun adapterScreenVertical(): Boolean {
        return true
    }

	private fun registerOnlineStatus() {
		vm.addRtcHandler(onlineStatusCallback)
	}

    private fun showLoading(content: String) {
        binding.loading.run {
            setCardColor(Color.TRANSPARENT)
            setFontColor(Color.parseColor("#FFFFFF"))
            showLoading(content)
        }
    }

    override fun onDestroy() {
        vm.removeRtcHandler(onlineStatusCallback)
        vm.logoutRtm()
        super.onDestroy()
    }
}
