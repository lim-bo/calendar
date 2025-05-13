#include "event_entry.h"
#include "ui_event_entry.h"

event_entry::event_entry(EventData data, QString viewerUID,  QWidget *parent)
    : QWidget(parent)
    , ui(new Ui::event_entry), data(data), viewerUID(viewerUID)
{
    ui->setupUi(this);

    ui->name->setText(data.name);
    ui->type->setText(data.type);
    ui->prior->setText(QString::number(data.prior));
    ui->vremyas->setText(data.start.toString() + "   " + data.end.toString());
    ui->desc->setText(data.desc);
}

const EventData event_entry::getData() const {
    return this->data;
}

event_entry::~event_entry()
{
    delete ui;
}

void event_entry::on_pushButton_clicked()
{
    emit deleted(this);
}


void event_entry::on_pushButton_3_clicked()
{
    ChatWindow chat(data.ID, viewerUID);
    chat.exec();
}

