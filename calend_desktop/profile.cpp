#include "profile.h"
#include "ui_profile.h"

profile::profile(QString uid, QWidget *parent)
    : QDialog(parent)
    , ui(new Ui::profile), cli(settings.value("host").toString(), settings.value("port").toString()), uid(uid)
{
    ui->setupUi(this);
    loadUserData();
}

profile::~profile()
{
    delete ui;
}

void profile::on_pushButton_clicked()
{
    ui->pushButton->setDisabled(true);
    credentials_reg data;
    data.mail = ui->email->text();
    data.pos = ui->position->text();
    data.department = ui->depart->text();
    QStringList fioParts = ui->fio->text().split(" ");
    if (fioParts.size() >= 3) {
        data.f_name = fioParts[0];
        data.s_name = fioParts[1];
        data.t_name = fioParts[2];
    } else {

        ui->result->setText("Введите ФИО полностью");
        return;
    }
    bool ok = cli.update(data, uid);
    if (!ok) {
        ui->result->setText("Ошибка");
        ui->pushButton->setEnabled(true);
    }
    else{
        this->close();
    }

}

void profile::loadUserData()
{

    credentials_reg pro = cli.getUserData(uid);
    QString fio = pro.f_name + " " + pro.s_name + " " + pro.t_name;
    ui->fio->setText(fio);
    ui->email->setText(pro.mail);
    ui->depart->setText(pro.department);
    ui->position->setText(pro.pos);



}
