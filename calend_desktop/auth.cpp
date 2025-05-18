#include "auth.h"
#include "ui_auth.h"

auth::auth(bool* fl, QString* uid, QWidget *parent)
    : QDialog(parent)
    , ui(new Ui::auth), fl(fl), uid(uid), cli(settings.value("host").toString(), settings.value("port").toString())
{
    ui->setupUi(this);
    ui->reg_button->setDisabled(true);
    QSettings settings;
}

auth::~auth()
{
    delete ui;
}

void auth::on_login_button_clicked()
{

    ui->login_button->setDisabled(true);
    credentials data;
    data.mail = ui->mail_edit->text();
    data.pass = ui->pass_edit->text();
    bool ok = cli.login(data, uid);
    if (!ok) {
        ui->result->setText("Ошибка");
        ui->login_button->setEnabled(true);
    } else {
        *fl = true;
        this->close();
    }
}
void auth::on_reg_button_clicked()
{
    ui->reg_button->setDisabled(true);
    credentials_reg data;
    data.mail = ui->mail_edit->text();
    data.pass = ui->pass_edit->text();

    QStringList fioParts = ui->fio->text().split(" ");
    if (fioParts.size() >= 3) {
        data.f_name = fioParts[0];
        data.s_name = fioParts[1];
        data.t_name = fioParts[2];
    } else {

        ui->result->setText("Введите ФИО полностью");
        return;
    }
    data.department = ui->depart->text();
    data.pos = ui->position->text();

    bool ok = cli.registration(data);
    if (!ok) {
        ui->result->setText("Ошибка регистрации");
        ui->reg_button->setEnabled(true);
    } else {
        ui->result->setText("Регистрация успешна");
        ui->reg_button->setEnabled(true);
        on_login_button_clicked();
    }

}

void auth::on_checkBox_toggled(bool checked)
{
    if (checked) {
        ui->box->setEnabled(true);
        ui->reg_button->setEnabled(true);
    }
    else {
        ui->box->setDisabled(true);
        ui->reg_button->setDisabled(true);
    }
}




