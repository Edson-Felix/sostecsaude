package com.wwmm.sostecsaude.ui.empresas

import android.os.Bundle
import androidx.fragment.app.Fragment
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup

import com.wwmm.sostecsaude.R
import kotlinx.android.synthetic.main.fragment_tab_defeito.*

class TabDefeito : Fragment() {

    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View? {
        return inflater.inflate(R.layout.fragment_tab_defeito, container, false)
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)

        arguments.apply {
            textView_defeito.text = this!!.getString("Defeito").toString()
            textView_quantidade.text = this.getString("Quantidade").toString()
        }
    }
}
