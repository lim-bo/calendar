#include "message_entry.h"
#include "ui_message_entry.h"

message_entry::message_entry(QWidget *parent)
    : QWidget(parent)
    , ui(new Ui::message_entry)
{
    ui->setupUi(this);
}

message_entry::~message_entry()
{
    delete ui;
}

void message_entry::setAttributes(QString sender, QString content) {
    ui->contentLabel->setText(content);
    ui->senderLabel->setText(sender);
}
