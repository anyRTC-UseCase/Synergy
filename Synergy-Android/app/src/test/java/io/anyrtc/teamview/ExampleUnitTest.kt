package io.anyrtc.teamview

import org.junit.Test

import org.junit.Assert.*

/**
 * Example local unit test, which will execute on the development machine (host).
 *
 * See [testing documentation](http://d.android.com/tools/testing).
 */
class ExampleUnitTest {
    @Test
    fun addition_isCorrect() {
        //assertEquals(4, 2 + 2)
        println((0 until 9).fold("") { acc, _ ->
            "$acc${(0 .. 9).random()}"
        })
        println((0 until 9).fold("") { acc, _ ->
            "$acc${(0 .. 9).random()}"
        })
        println((0 until 9).fold("") { acc, _ ->
            "$acc${(0 .. 9).random()}"
        })
        println((0 until 9).fold("") { acc, _ ->
            "$acc${(0 .. 9).random()}"
        })
    }
}