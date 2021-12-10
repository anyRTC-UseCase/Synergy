package io.anyrtc.teamview.activity

import android.graphics.Color
import android.os.Build
import android.os.Bundle
import android.widget.TextView
import android.widget.Toast
import androidx.core.content.ContextCompat
import androidx.recyclerview.widget.LinearLayoutManager
import android.view.View
import android.view.Window
import android.view.WindowInsets
import android.view.WindowManager
import androidx.core.view.GravityCompat
import androidx.activity.viewModels
import androidx.recyclerview.widget.RecyclerView
import io.anyrtc.teamview.R
import io.anyrtc.teamview.databinding.ActivityMainBinding
import io.anyrtc.teamview.utils.*
import io.anyrtc.teamview.vm.MainVM
import io.anyrtc.teamview.utils.RtcManager
import org.ar.rtm.LocalInvitation
import java.util.*
import kotlin.collections.HashSet

class MainActivity : BaseActivity() {

    private lateinit var binding: ActivityMainBinding
    private val vm: MainVM by viewModels()

    private lateinit var uid: str
    private var roomId = ""
    private var roomName = ""
    private var roomTs = 0
    private var userName = ""
	private var notMore = false
	private var loadingMore = false
	private var isRefreshing = false

    private val itemBtnStatusTextArr = arrayOf("邀请中", "通话中", "邀请", "离线")
    private val itemStatusColorArr = arrayOf(
        R.color.blue, R.color.red, R.color.blue, R.color.gray
    )

    private val statusToWeightsArr = arrayOf(1, 2, 0, 3)

    private val adapter by lazy {
        Adapter(arrayListOf(), this@MainActivity::bindItem, R.layout.item_main)
    }
	private val mLocalDataMap = HashMap<String, ItemData>()
    private var mTimer: Timer? = null
    private val mRtcHandler = RtcHandler()

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        requestWindowFeature(Window.FEATURE_NO_TITLE)
        binding = ActivityMainBinding.inflate(layoutInflater)
        hideBottomNavigationBar()
        setContentView(binding.root)

        // full screen, hide status bar
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.R)
            window.insetsController?.hide(WindowInsets.Type.statusBars())
        else this.window.setFlags(
            WindowManager.LayoutParams.FLAG_FULLSCREEN,
            WindowManager.LayoutParams.FLAG_FULLSCREEN
        )

        val userTs = intent.getIntExtra("userTs", 0)
        val userType = intent.getIntExtra("userType", 0)
        val workName = intent.getStringExtra("workName")!!
        val rtcToken = intent.getStringExtra("rtcToken")!!
        roomId = intent.getStringExtra("roomId")!!
        roomName = intent.getStringExtra("roomName")!!
        roomTs = intent.getIntExtra("roomTs", 0)
        userName = intent.getStringExtra("userName")!!
        uid = SpUtil.get().getString(MyConstants.UID, "")!!

        if (uid.isEmpty()) {
            // error
        }
		vm.addRtcHandler(mRtcHandler)

        initRtc(rtcToken, roomId)
        initWidget()
        initData()
    }

    private fun initWidget() {
        binding.recycle.run {
            layoutManager = LinearLayoutManager(
                this@MainActivity, LinearLayoutManager.VERTICAL, false
            )
            adapter = this@MainActivity.adapter
            addItemDecoration(MainItemDecoration())

			addOnScrollListener(object : RecyclerView.OnScrollListener() {
				override fun onScrollStateChanged(view: RecyclerView, newState: Int) {
					if (!canScrollVertically(1)) {
						loadOrRefreshSpecialList()
					}
				}
			})
        }
        binding.hangUp.setOnClickListener { // TODO: may need a dialog to confirm
            finish()
        }
        binding.invite.setOnClickListener {
            binding.drawer.openDrawer(GravityCompat.END)
        }
        binding.close.setOnClickListener {
            binding.drawer.closeDrawer(GravityCompat.END)
        }
        binding.refresh.setOnRefreshListener {
			notMore = false
			loadOrRefreshSpecialList(false)
        }
    }

    private fun initData() {
        vm.specialListResponse.observe(this) {
			binding.refresh.isRefreshing = false
            if (it.data == null) {
                Toast.makeText(this@MainActivity, it.errMsg, Toast.LENGTH_SHORT).show()
				loadingMore = false
                return@observe
            }
			// 1. Add if local hashMap don't have that uid
			// 2. Also add it to the new ArrayList
			// 3. Sorted new ArrayList & add to adapter & notify refresh

            val list = it.data.list
			if (list.size == adapter.data.size && !isRefreshing) {
				binding.refresh.isRefreshing = false
				notMore = true
				loadingMore = false
                return@observe
            }
			isRefreshing = false

			val copiedLocalMap = HashMap(mLocalDataMap)
			val newList = mutableListOf<ItemData>()
			var index = 0
			val adapterSize = adapter.itemCount

			val sortedList = list.map { item ->
				// 0=邀请中，1=通话中，2=空闲，3=离线
				val status =  if (item.userState == 1) if (item.roomId != roomId) 2 else 1 else item.userState
				val itemData = ItemData(
					item.uid,
					item.userName,
					status,
					weights = statusToWeightsArr[status]
				)
				// check status change if not null
				if (mLocalDataMap[item.uid] != null) {
					copiedLocalMap.remove(item.uid)
					val localItemData = mLocalDataMap[item.uid]!!
					if (localItemData.status == 0 && itemData.status == 2) {
						itemData.status = 0
						itemData.weights = statusToWeightsArr[0]
					} else {
						mLocalDataMap.remove(item.uid)
					}
				}

				if (index++ >= adapterSize) {
					newList.add(itemData)
				}

				itemData
			}.sortedBy { item -> item.weights }

			if (copiedLocalMap.isNotEmpty())
			mLocalDataMap.removeAll(copiedLocalMap)

			adapter.data.addAll(newList)
			adapter.notifyItemRangeInserted(adapterSize, newList.size)

			adapter.data.clear()
			adapter.data.addAll(sortedList)
			adapter.notifyItemRangeChanged(0, sortedList.size)

			binding.title.text = String.format("邀请协调（%d）", adapter.data.size)
			loadingMore = false
        }

		loadOrRefreshSpecialList()
    }

	private fun <K, V> HashMap<K, V>.removeAll(m: Map<in K, V>): Boolean {
		var modified = false
		if (size > m.size) {
			val i = m.iterator()
			while (i.hasNext()) {
				modified = modified or (remove(i.next().key) != null)
			}
		} else {
			val i = iterator()
			while (i.hasNext()) {
				if (m.containsKey(i.next().key)) {
					i.remove()
					modified = true
				}
			}
		}
		return modified
	}

    private fun initRtc(rtcToken: str, roomId: str) {
        vm.sdkJoinRoom(this, rtcToken, roomId, uid) { joinFailed ->
            if (joinFailed) {
                Toast.makeText(binding.root.context, "join rtc failed", Toast.LENGTH_SHORT).show()
                throw IllegalStateException("join rtc failed.")
            }
        }
        val textureView = vm.createRendererView(this, uid)
        binding.cameraParent.run {
            removeAllViews()
            addView(textureView)
            vm.switchCamera()
        }
    }

	private fun loadOrRefreshSpecialList(isLoadMore: Boolean = true) {
		if (isLoadMore && notMore)
		return
		if (loadingMore)
		return

		loadingMore = true
		binding.refresh.isRefreshing = true
		if (isLoadMore) {
			vm.getSpecialList(1, adapter.data.size + 10)
			return
		}

		// 下拉刷新清除本地数据
		binding.title.text = "邀请协调"
		val removeSize = adapter.data.size
		adapter.data.clear()
		adapter.notifyItemRangeRemoved(0, removeSize)
		vm.getSpecialList(1, 10)
	}

    override fun adapterScreenVertical(): Boolean {
        return true
    }

    private fun bindItem(holder: Adapter.Holder, data: ItemData, position: int) {
        val nickname = holder.findView<TextView>(R.id.nickname)
        val statusView = holder.findView<TextView>(R.id.invite_status)

        nickname.text = data.nickname
        statusView.text = itemBtnStatusTextArr[data.status]

        statusView.visibility = if (data.status == 3) View.GONE else View.VISIBLE
        statusView.setBackgroundResource(
            if (data.status == 2) R.drawable.shape_item_req_btn else android.R.color.transparent
        )
        statusView.setTextColor(ContextCompat.getColor(this, itemStatusColorArr[data.status]))

        statusView.setOnClickListener {
            if (data.status != 2 || binding.refresh.isRefreshing)
                return@setOnClickListener

            val sendCallBlock = {
                vm.sendCall(
                    data.uid,
                    roomId,
                    roomName,
                    roomTs,
                    userName
                ) { success, err, localInvitation ->
                    if (!success) {
                        if (null != err)
                            Toast.makeText(this@MainActivity, err.toString(), Toast.LENGTH_SHORT).show()
                        return@sendCall
                    }
                    data.status = 0
                    data.weights = statusToWeightsArr[data.status]
                    data.sentInvitation = localInvitation
                    adapter.notifyItemChanged(position)
					mLocalDataMap[data.uid] = data
                }
            }
            if (data.sentInvitation != null) {
                vm.cancelCall(data.sentInvitation!!) { cancelSuccess, cancelErr ->
                    if (!cancelSuccess) {
                        if (null != cancelErr)
                            Toast.makeText(this, cancelErr.toString(), Toast.LENGTH_LONG).show()
                        return@cancelCall
                    }
                    sendCallBlock.invoke()
                }
            } else {
                sendCallBlock.invoke()
            }
        }
    }

    private inner class RtcHandler : RtcManager.RtcHandler() {
        private var waitingDestroyThread: Thread? = null
        override fun onCallLocalAccept(var1: LocalInvitation?) {
            var1 ?: return
            var i = 0
            while (1 < adapter.data.size) {
                if (adapter.data[i].uid == var1.calleeId) {
                    break
                }
                i++
            }
            if (i == adapter.data.size) {
                // not found.
                return
            }
            runOnUiThread {
                val nickname = adapter.data[i].nickname
                binding.notify.showMessage("$nickname 同意了你的邀请")
				refreshSpecialList()
            }
        }

        override fun onCallLocalReject(var1: LocalInvitation?) {
            var1 ?: return
            var i = 0
            while (1 < adapter.data.size) {
                if (adapter.data[i].uid == var1.calleeId) {
                    break
                }
                i++
            }
            if (i == adapter.data.size) {
                // not found.
                return
            }
			mLocalDataMap[adapter.data[i].uid]?.status = 2
            runOnUiThread {
                val nickname = adapter.data[i].nickname
                binding.notify.showMessage("$nickname 拒绝了你的邀请")
				refreshSpecialList()
            }
        }

        override fun onCallFailure(var1: LocalInvitation?, var2: Int) {
            var1 ?: return
            if (var2 == 0) return

            var i = 0
            while (true) {
                if (adapter.data[i].uid == var1.calleeId)
                    break
                i++
            }
            if (i == adapter.data.size) {
                // not found.
                return
            }
			mLocalDataMap[adapter.data[i].uid]?.status = 2
            runOnUiThread {
                Toast.makeText(applicationContext, "对方60秒内未接受邀请", Toast.LENGTH_SHORT).show()
				refreshSpecialList()
            }
        }

        override fun onUserJoined(uid: String?, elapsed: Int) {
            uid ?: return
            var i = 0
            while (i < adapter.data.size) {
                val item = adapter.data[i]
                if (item.uid == uid)
                break

                ++i
            }
            if (i < adapter.data.size) runOnUiThread {
                val nickname = adapter.data[i].nickname
                binding.notify.showMessage("$nickname 进入了协同")
				refreshSpecialList()
            }
        }

        override fun onUserLeave(uid: String?, reason: Int) {
            uid ?: return
            var i = 0
            while (i < adapter.data.size) {
                val item = adapter.data[i]
                if (item.uid == uid)
                break

                ++i
            }
            if (i < adapter.data.size) runOnUiThread {
                val nickname = adapter.data[i].nickname
                binding.notify.showMessage("$nickname 离开了协同")
				refreshSpecialList()
            }
        }

        override fun onRtcTokenExpired() {
            runOnUiThread {
                Toast.makeText(this@MainActivity, "体验时间已到", Toast.LENGTH_LONG).show()
                finish()
            }
        }

        override fun selfNetworkChange(disconnect: bool) {
            if (!disconnect && waitingDestroyThread != null) {
                waitingDestroyThread?.interrupt()
                waitingDestroyThread = null
                runOnUiThread { binding.loading.hideLoading() }
            } else if (disconnect && waitingDestroyThread == null) {
                runOnUiThread {
                    binding.loading.run {
                        setCardColor(Color.TRANSPARENT)
                        setFontColor(Color.WHITE)
                        showLoading("连网中")
                    }
                }
                waitingDestroyThread = ShutdownThread {
                    runOnUiThread {
                        Toast.makeText(this@MainActivity, "连接已断开", Toast.LENGTH_LONG).show()
                        finish()
                    }
                }
                waitingDestroyThread?.start()
            }
        }

        override fun onOtherClientLogin() {
            runOnUiThread {
                finish()
            }
        }

		private fun refreshSpecialList() {
			binding.refresh.isRefreshing = true
			isRefreshing = true
			vm.getSpecialList(1, adapter.itemCount)
		}
    }

    private fun hideBottomNavigationBar() {
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.R)
            window.insetsController?.hide(WindowInsets.Type.navigationBars())
        else
            window.decorView.systemUiVisibility =
                window.decorView.systemUiVisibility or View.SYSTEM_UI_FLAG_HIDE_NAVIGATION or View.SYSTEM_UI_FLAG_LOW_PROFILE
    }

    override fun onDestroy() {
        if (roomId.isNotBlank())
            vm.leaveRoom(roomId, uid)
        vm.releaseRtc()
		vm.removeRtcHandler(mRtcHandler)
        mTimer?.let {
            it.cancel()
            it.purge()
        }
        mTimer = null

		adapter.data.forEach { item ->
			if (item.sentInvitation != null) {
                vm.cancelCall(item.sentInvitation!!)
            }
		}
        super.onDestroy()
    }

    data class ItemData(
        val uid: str,
        val nickname: str,
        /*
         * 0=邀请中，1=通话中，2=空闲，3=离线
         * 0=邀请中，1=在线，2=空闲，3=离线, 4=通话中
         */
        var status: int,
        var weights: int,
        var sentInvitation: LocalInvitation? = null
    )

    private class ShutdownThread(private val destroyCallback: () -> Unit) : Thread() {
        private var lifetime = 60
        override fun run() {
            while (!isInterrupted) {
                try {
                    sleep(1000)
                    if (--lifetime <= 0) {
                        destroyCallback.invoke()
                        break
                    }
                } catch (e: InterruptedException) {
                    break
                }
            }
        }
    }
}
