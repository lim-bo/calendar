#include "attachments_comp.h"
#include "ui_attachments_comp.h"

attachments_comp::attachments_comp(QWidget *parent)
    : QWidget(parent)
    , ui(new Ui::attachments_comp)
{
    ui->setupUi(this);
}

attachments_comp::~attachments_comp()
{
    delete ui;
}

void attachments_comp::setData(const QString &name, QByteArray data) {
    ui->labelName->setText(name);
    this->data = data;
}

void attachments_comp::on_pushButton_clicked()
{
    QString savePath = QStandardPaths::writableLocation(QStandardPaths::DownloadLocation);
    QSaveFile file(savePath + "/"+ui->labelName->text());
    if (file.open(QIODevice::WriteOnly)) {
        file.write(data);
        if (!file.commit())
            qDebug() << "Ошибка сохранения файла: " <<file.errorString();
        else
            QMessageBox::information(this, "Сохранено", "Файл сохранён в Загрузки");
    } else {
        qDebug() << "Ошибка открытия файла: " << file.errorString();
    }

}

