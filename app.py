# app.py
import streamlit as st
import pandas as pd
from datetime import datetime
from codelight import (
    load_data,
    recommend_bundles,
    get_available_slots,
    book_installation,
    vdsl_speed_from_distance,
    normalize_currency  # Yeni eklenen fonksiyon
)

# Sayfa ayarları
st.set_page_config(
    page_title="Turkcell Ev+Mobil Paket Danışmanı",
    page_icon="📱",
    layout="wide",
    initial_sidebar_state="expanded"
)

# Tema ve erişilebilirlik için CSS
css = """
<style>
    .st-emotion-cache-1y4p8pa {
        padding: 2rem 1rem 10rem;
    }
    .stButton>button {
        background-color: #4CAF50;
        color: white;
    }
    .metric-card {
        border-radius: 0.5rem;
        padding: 1rem;
        background-color: #f0f2f6;
        margin-bottom: 1rem;
    }
    /* Karanlık tema */
    .dark-theme {
        background-color: #1e1e1e;
        color: white;
    }
    .dark-theme .st-bq {
        background-color: #2d2d2d;
    }
    /* Yüksek kontrast modu */
    .high-contrast {
        filter: contrast(1.5);
    }
    /* Büyük font */
    .large-font {
        font-size: 1.2em;
    }
    /* Odaklanma göstergesi */
    *:focus {
        outline: 3px solid #4CAF50;
    }
</style>
"""

st.markdown(css, unsafe_allow_html=True)

# Veri yükleme
@st.cache_resource
def load_app_data():
    try:
        return load_data()
    except Exception as e:
        st.error(f"Veri yüklenirken hata oluştu: {str(e)}")
        st.stop()

data = load_app_data()

# ==================== SIDEBAR ====================
with st.sidebar:
    st.image("https://www.turkcell.com.tr/Content/images/logo.png", width=150)
    st.title("Paket Danışmanı")
    
    # Kullanıcı seçimi
    user_id = st.selectbox(
        "Kullanıcı Seçin",
        options=data['users']['user_id'].unique(),
        index=0
    )
    
    # Adres bilgisi
    user_data = data['users'][data['users']['user_id'] == user_id].iloc[0]
    address_id = user_data['address_id']
    st.write(f"**Adres ID:** {address_id}")
    
    # Teknoloji tercihi
    tech_options = ['fiber', 'vdsl', 'fwa']
    prefer_techs = st.multiselect(
        "Tercih Edilen Teknolojiler (Sıralı)",
        options=tech_options,
        default=tech_options
    )
    
    # Para birimi seçimi
    currency = st.selectbox(
        "Para Birimi",
        options=['TL', 'USD', 'EUR'],
        index=0
    )
    
    # Tema seçimi
    st.subheader("Tema Seçimi")
    theme = st.selectbox(
        "Tema",
        options=["Aydınlık", "Karanlık"]
    )
    
    # Erişilebilirlik
    st.subheader("Erişilebilirlik")
    high_contrast = st.checkbox("Yüksek Kontrast")
    large_font = st.checkbox("Büyük Font")
    
    # Hesapla butonu
    if st.button("Paketleri Hesapla", type="primary", use_container_width=True):
        st.session_state['calculate'] = True
    else:
        st.session_state['calculate'] = False

# Tema ve erişilebilirlik uygula
if theme == "Karanlık":
    st.markdown('<div class="dark-theme"></div>', unsafe_allow_html=True)

if high_contrast:
    st.markdown('<div class="high-contrast"></div>', unsafe_allow_html=True)
    
if large_font:
    st.markdown('<div class="large-font"></div>', unsafe_allow_html=True)

# ==================== ANA SAYFA ====================
st.title("Turkcell Ev+Mobil Paket Danışmanı")
st.markdown("""
Hane halkınız için en uygun paket kombinasyonlarını bulalım. 
Adresinizdeki altyapı durumuna ve kullanım ihtiyaçlarınıza göre öneriler sunuyoruz.
""")

# Kapsama bilgisi
st.subheader("📍 Adres Kapsama Durumu")
coverage = data['coverage'][data['coverage']['address_id'] == address_id].iloc[0]

cols = st.columns(3)
with cols[0]:
    st.metric("Fiber Altyapı", "✅ Uygun" if coverage['fiber'] == 1 else "❌ Yok")
with cols[1]:
    if coverage['vdsl'] == 1:
        speed = vdsl_speed_from_distance(coverage.get('distance_km', 0))
        st.metric("VDSL Altyapı", f"✅ Uygun (~{speed:.1f} Mbps)")
    else:
        st.metric("VDSL Altyapı", "❌ Yok")
with cols[2]:
    st.metric("FWA (4.5G)", "✅ Uygun" if coverage['fwa'] == 1 else "❌ Yok")

# Paket önerileri
if st.session_state.get('calculate', False):
    try:
        with st.spinner("En uygun paketler hesaplanıyor..."):
            recommendations, pred_data = recommend_bundles(user_id, data, prefer_techs)
            
            # Tahmin bilgileri
            st.subheader("📊 Tahmini Kullanım")
            pred_cols = st.columns(3)
            pred_cols[0].metric("Toplam Data", f"{pred_data['total_gb']:.1f} GB")
            pred_cols[1].metric("Toplam Dakika", f"{pred_data['total_min']:.0f} dk")
            pred_cols[2].metric("TV İzleme", f"{pred_data['total_tv_hours']:.1f} saat")
            
            # Öneriler
            st.subheader("🏆 Size Özel Paket Önerileri")
            
            for i, combo in enumerate(recommendations):
                # Para birimine göre dönüştür
                total_cost = normalize_currency(combo['monthly_total'], 'TL', currency)
                savings = normalize_currency(combo['savings'], 'TL', currency)
                
                with st.expander(f"{i+1}. {combo['combo_label']} | {total_cost:.2f} {currency}", expanded=i==0):
                    # Üst bilgi
                    cols = st.columns([1, 2])
                    
                    with cols[0]:
                        st.metric("Aylık Toplam", f"{total_cost:.2f} {currency}")
                        
                        if savings > 0:
                            st.metric("Tahmini Tasarruf", 
                                     f"{savings:.2f} {currency}", 
                                     delta=f"%{(savings/total_cost*100):.1f}")
                        else:
                            st.metric("Maliyet Farkı", 
                                     f"{-savings:.2f} {currency}", 
                                     delta_color="inverse")
                        
                        st.caption(combo['reasoning'])
                    
                    with cols[1]:
                        # Mobil detayları
                        st.markdown("**📱 Mobil Hatlar**")
                        for line in combo['items']['mobile']:
                            cost = normalize_currency(line['cost'], 'TL', currency)
                            st.markdown(f"""
                            - **{line['line_id']}**: {line['plan_name']}  
                              Kota: {line['quota_gb']}GB + {line['quota_min']}dk  
                              Aşım: {line['overage_gb']}TL/GB, {line['overage_min']}TL/dk  
                              **Maliyet**: {cost:.2f} {currency}
                            """)
                        
                        # Ev interneti
                        if 'home' in combo['items']:
                            home = combo['items']['home']
                            home_cost = normalize_currency(home['monthly_price'], 'TL', currency)
                            st.markdown(f"""
                            **🏠 Ev İnterneti**  
                            - {home['name']} ({home['tech'].upper()})  
                            Hız: {home['down_mbps']} Mbps  
                            **Maliyet**: {home_cost:.2f} {currency}
                            """)
                        
                        # TV
                        if 'tv' in combo['items'] and combo['items']['tv']:
                            tv = combo['items']['tv']
                            tv_cost = normalize_currency(tv['monthly_price'], 'TL', currency)
                            st.markdown(f"""
                            **📺 TV Paketi**  
                            - {tv['name']}  
                            Dahil HD Saat: {tv['hd_hours_included']}  
                            **Maliyet**: {tv_cost:.2f} {currency}
                            """)
                    
                    # Randevu butonu (sadece ev interneti içeren paketlerde)
                    if 'home' in combo['items']:
                        if st.button(f"{combo['combo_label']} Paketini Seç", key=f"select_{i}"):
                            st.session_state['selected_combo'] = combo
                            st.session_state['show_slots'] = True
                            st.rerun()
            
            # Randevu seçim ekranı
            if st.session_state.get('show_slots', False) and 'selected_combo' in st.session_state:
                combo = st.session_state['selected_combo']
                tech = combo['tech_label']
                
                st.subheader(f"⏰ Kurulum Randevusu - {combo['combo_label']}")
                st.info(f"{tech.upper()} altyapısı için uygun kurulum slotları:")
                
                slots = get_available_slots(address_id, tech, data)
                if slots:
                    slot_options = [f"{s['slot_id']} | {s['slot_start'].strftime('%d.%m.%Y %H:%M')} - {s['slot_end'].strftime('%H:%M')}" 
                                   for s in slots]
                    selected_slot = st.selectbox("Randevu Slotu Seçin", slot_options)
                    
                    if st.button("Randevuyu Onayla", type="primary"):
                        slot_id = selected_slot.split(' | ')[0]
                        if book_installation(user_id, slot_id, data):
                            st.success("Randevunuz başarıyla alındı!")
                            
                            # Hat taşı/ek hat ekle akışı
                            st.subheader("Hat Taşıma / Ek Hat Ekleme")
                            st.info("Mevcut hatlarınızı yeni pakete taşımak veya yeni hat eklemek için aşağıdaki adımları takip edin:")
                            
                            hat_tasima = st.checkbox("Mevcut hatlarımı yeni pakete taşımak istiyorum")
                            ek_hat = st.checkbox("Yeni hat eklemek istiyorum")
                            
                            if hat_tasima or ek_hat:
                                st.text_input("Hat numarası (ör: 5XXXXXXXXX)")
                                st.selectbox("İşlem türü", ["Hat taşıma", "Yeni hat"])
                                st.button("Devam Et")
                            
                            # Modem teslimatı akışı
                            st.subheader("Modem Teslimatı")
                            st.info("Modem teslimat seçenekleri:")
                            
                            modem_secenek = st.radio(
                                "Modem teslimat şekli",
                                ["Adrese teslim", "Turkcell mağazasından teslim al"]
                            )
                            
                            if modem_secenek == "Adrese teslim":
                                st.text_input("Teslimat adresi")
                                st.date_input("Teslimat tarihi")
                            else:
                                st.selectbox("Mağaza seçin", ["Kadıköy", "Levent", "İstinye Park"])
                            
                            if st.button("Siparişi Tamamla", type="primary"):
                                st.success("Siparişiniz başarıyla oluşturuldu!")
                                st.balloons()
                                st.session_state['show_slots'] = False
                        else:
                            st.error("Randevu alınamadı, lütfen tekrar deneyin.")
                else:
                    st.warning("Bu teknoloji için uygun randevu slotu bulunamadı.")
                    if st.button("Geri Dön"):
                        st.session_state['show_slots'] = False
                        st.rerun()
        
    except Exception as e:
        st.error(f"Hata oluştu: {str(e)}")
        st.stop()

# ==================== FOOTER ====================
st.markdown("---")
st.caption("""
Turkcell Ev+Mobil Paket Danışmanı - Tüm hakları saklıdır.  
Bu bir demo uygulamasıdır, gerçek paketler ve fiyatlar için Turkcell'i ziyaret edin.
""")