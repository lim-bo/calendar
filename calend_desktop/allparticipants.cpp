#include "allparticipants.h"
#include "ui_allparticipants.h"

allparticipants::allparticipants(const EventData &eventData, QWidget *parent)
    : QDialog(parent)
    , ui(new Ui::allparticipants), eventData(eventData)
{
    ui->setupUi(this);

    Client cli(settings.value("host").toString(), settings.value("port").toString());
    QString participantsText;
    QVector<Participant> parts = cli.getParticipants(eventData.ID);
    for (const Participant &participant : parts) {
        QString status = participant.accepted ? "Подтвердил" : "Не подтвердил";
        participantsText += QString("%1 - %2\n").arg(participant.uid).arg(status);
    }
    ui->textEdit->setPlainText(participantsText);
    ui->textEdit->setReadOnly(true);
}

allparticipants::~allparticipants()
{
    delete ui;
}
